import re
import json
import hashlib
import uuid
from typing import Dict, Any, Optional, List, Literal
from datetime import datetime

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.semantic_conventions.resource import ResourceAttributes
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import Response
import httpx

from .config import AgentConfig


class SecurityContext:
    """Security context for request tracking"""
    def __init__(self, request_id: str, trace_id: str, session_id: Optional[str] = None,
                 user_id: Optional[str] = None, taint_tags: Optional[List[str]] = None):
        self.request_id = request_id
        self.trace_id = trace_id
        self.session_id = session_id
        self.user_id = user_id
        self.taint_tags = taint_tags or []


class SecurityAgent:
    """Security Platform Agent for Python applications"""

    def __init__(self, config: AgentConfig):
        self.config = config
        self.redaction_rules = [re.compile(rule, re.IGNORECASE) for rule in config.redaction_rules]
        self.local_policy: Dict[str, Literal["block", "observe"]] = {}
        self.tracer_provider = None
        self.tracer = None
        
        self._setup_telemetry()
        self._load_local_policy()

    def _setup_telemetry(self):
        """Initialize OpenTelemetry tracing"""
        resource = Resource.create({
            ResourceAttributes.SERVICE_NAME: self.config.service_name,
            ResourceAttributes.SERVICE_VERSION: "0.1.0",
            ResourceAttributes.DEPLOYMENT_ENVIRONMENT: self.config.environment,
        })

        exporter = OTLPSpanExporter(
            endpoint=self.config.otlp_endpoint,
            headers={"Authorization": f"Bearer {self.config.auth_key}"}
        )

        self.tracer_provider = TracerProvider(resource=resource)
        self.tracer_provider.add_span_processor(
            BatchSpanProcessor(exporter, max_queue_size=self.config.telemetry_batch_size)
        )
        trace.set_tracer_provider(self.tracer_provider)
        self.tracer = trace.get_tracer(self.config.service_name)

    async def _load_local_policy(self):
        """Load security policy from control plane"""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"{self.config.control_plane_url}/api/v1/policy",
                    headers={"Authorization": f"Bearer {self.config.auth_key}"},
                    timeout=5.0
                )
                if response.status_code == 200:
                    policy = response.json()
                    self.local_policy = policy.get("rules", {})
        except Exception as e:
            print(f"Warning: Failed to load policy from control plane: {e}")

    def _redact_data(self, data: Any) -> Any:
        """Redact sensitive data based on rules"""
        if isinstance(data, str):
            redacted = data
            for rule in self.redaction_rules:
                redacted = rule.sub("[REDACTED]", redacted)
            return redacted

        if isinstance(data, dict):
            return {k: self._redact_data(v) if not any(r.search(k) for r in self.redaction_rules) 
                    else "[REDACTED]" for k, v in data.items()}

        if isinstance(data, list):
            return [self._redact_data(item) for item in data]

        return data

    def _check_policy(self, context: SecurityContext) -> Literal["block", "observe"]:
        """Check if request should be blocked based on policy"""
        for tag in context.taint_tags:
            action = self.local_policy.get(tag)
            if action == "block":
                return "block"

        return "block" if self.config.local_policy == "block" else "observe"

    def capture_event(self, name: str, data: Optional[Dict[str, Any]] = None,
                     taint_tags: Optional[List[str]] = None):
        """Capture a security event"""
        if not self.tracer:
            return

        span = self.tracer.start_span(name)
        
        if taint_tags:
            span.set_attribute("security.taint_tags", ",".join(taint_tags))

        if data:
            redacted_data = self._redact_data(data)
            for key, value in redacted_data.items():
                span.set_attribute(key, str(value))

        span.end()


class SecurityMiddleware(BaseHTTPMiddleware):
    """Starlette/FastAPI middleware for security observability"""

    def __init__(self, app, agent: SecurityAgent):
        super().__init__(app)
        self.agent = agent

    async def dispatch(self, request: Request, call_next):
        """Process request through security middleware"""
        tracer = self.agent.tracer
        if not tracer:
            return await call_next(request)

        request_id = str(uuid.uuid4())
        span = tracer.start_span("http_request")
        trace_id = format(span.get_span_context().trace_id, '032x')

        context = SecurityContext(
            request_id=request_id,
            trace_id=trace_id,
            session_id=request.headers.get("x-session-id"),
            user_id=request.headers.get("x-user-id"),
        )

        # Capture request metadata
        span.set_attribute("http.method", request.method)
        span.set_attribute("http.url", str(request.url))
        span.set_attribute("http.route", request.url.path)
        span.set_attribute("request.id", request_id)

        # Redact headers
        redacted_headers = self.agent._redact_data(dict(request.headers))
        span.set_attribute("http.request.headers", json.dumps(redacted_headers))

        # Check body if available
        if hasattr(request, "_body"):
            try:
                body = await request.json()
                redacted_body = self.agent._redact_data(body)
                span.set_attribute("http.request.body", json.dumps(redacted_body))
            except:
                pass

        # Check policy
        action = self.agent._check_policy(context)
        if action == "block":
            span.set_attribute("security.action", "blocked")
            span.end()
            return Response(
                content=json.dumps({"error": "Request blocked by security policy"}),
                status_code=403,
                media_type="application/json"
            )

        # Add context to request state
        request.state.security_context = context

        # Process request
        try:
            response = await call_next(request)
            span.set_attribute("http.status_code", response.status_code)
            span.set_attribute("security.action", action)
            return response
        except Exception as e:
            span.set_attribute("error", True)
            span.set_attribute("error.message", str(e))
            raise
        finally:
            span.end()


def create_agent(config: AgentConfig) -> SecurityAgent:
    """Create a new SecurityAgent instance"""
    return SecurityAgent(config)
