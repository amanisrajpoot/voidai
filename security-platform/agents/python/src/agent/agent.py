"""Core agent implementation."""

import os
import logging
from typing import Optional, Dict, Any, List
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.semantic_conventions.resource import ResourceAttributes

from .types import AgentConfig
from .redaction import Redactor
from .policy import PolicyEngine

logger = logging.getLogger(__name__)


class Agent:
    """Security Platform Agent for Python applications."""

    def __init__(self, config: AgentConfig):
        """Initialize the agent with configuration."""
        self.config = config
        self.tracer_provider: Optional[TracerProvider] = None
        self.redactor = Redactor(config.get("redaction_rules", []))
        self.policy_engine = PolicyEngine(config.get("policy", {}))

    def start(self) -> None:
        """Start the agent and initialize OpenTelemetry."""
        try:
            resource = Resource.create({
                ResourceAttributes.SERVICE_NAME: self.config.get(
                    "service_name", "python-service"
                ),
                ResourceAttributes.DEPLOYMENT_ENVIRONMENT: self.config.get(
                    "environment", "production"
                ),
            })

            otlp_endpoint = self.config.get(
                "otlp_endpoint",
                f"{self.config.get('control_plane_url', 'http://localhost:4318')}/v1/traces",
            )

            exporter = OTLPSpanExporter(
                endpoint=otlp_endpoint,
                headers=self._get_auth_headers(),
            )

            self.tracer_provider = TracerProvider(resource=resource)
            span_processor = BatchSpanProcessor(
                exporter,
                max_queue_size=self.config.get("telemetry", {}).get("max_queue_size", 2048),
                max_export_batch_size=self.config.get("telemetry", {}).get("batch_size", 100),
            )
            self.tracer_provider.add_span_processor(span_processor)
            trace.set_tracer_provider(self.tracer_provider)

            logger.info("Security Platform Agent started")
        except Exception as e:
            logger.error(f"Failed to start agent: {e}", exc_info=True)
            raise

    def stop(self) -> None:
        """Stop the agent and shutdown OpenTelemetry."""
        if self.tracer_provider:
            self.tracer_provider.shutdown()
            self.tracer_provider = None
            logger.info("Security Platform Agent stopped")

    def redact(self, data: Any) -> Any:
        """Redact PII from data."""
        return self.redactor.redact(data)

    def check_policy(
        self, action: str, context: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Check if an action is allowed by policy."""
        return self.policy_engine.check(action, context)

    def _get_auth_headers(self) -> Dict[str, str]:
        """Get authentication headers for OTLP exporter."""
        headers = {}
        auth_key = self.config.get("auth_key") or os.getenv(
            "SECURITY_PLATFORM_AUTH_KEY"
        )
        if auth_key:
            headers["Authorization"] = f"Bearer {auth_key}"
        return headers
