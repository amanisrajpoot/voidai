"""Main agent implementation"""

import re
import json
import hashlib
from typing import Any, Dict, List, Optional
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.semantic_conventions.resource import ResourceAttributes

from .config import AgentConfig


class SecurityAgent:
    """Security observability agent"""

    def __init__(self, config: AgentConfig):
        self.config = config
        self.tracer_provider = None
        self.event_queue: List[Dict[str, Any]] = []
        self._initialize_telemetry()

    def _initialize_telemetry(self):
        """Initialize OpenTelemetry tracing"""
        resource = Resource.create({
            ResourceAttributes.SERVICE_NAME: self.config.service_name,
            ResourceAttributes.DEPLOYMENT_ENVIRONMENT: self.config.env,
        })

        self.tracer_provider = TracerProvider(resource=resource)

        exporter = OTLPSpanExporter(
            endpoint=f"{self.config.otlp_endpoint}/v1/traces",
            headers={"Authorization": f"Bearer {self.config.auth_key}"} if self.config.auth_key else None,
        )

        self.tracer_provider.add_span_processor(BatchSpanProcessor(exporter))
        trace.set_tracer_provider(self.tracer_provider)

    def redact_data(self, data: Any) -> Any:
        """Redact sensitive data based on configuration rules"""
        if isinstance(data, dict):
            redacted = {}
            for key, value in data.items():
                if self._should_redact(key):
                    redacted[key] = "[REDACTED]"
                else:
                    redacted[key] = self.redact_data(value)
            return redacted
        elif isinstance(data, list):
            return [self.redact_data(item) for item in data]
        elif isinstance(data, str):
            # Hash PII strings if they match patterns
            if self._is_pii(data):
                return self._hash_pii(data)
            return data
        else:
            return data

    def _should_redact(self, key: str) -> bool:
        """Check if a key should be redacted"""
        key_lower = key.lower()
        for rule in self.config.redaction_rules:
            pattern = rule.replace("*", ".*")
            if re.match(pattern, key_lower, re.IGNORECASE):
                return True
        return False

    def _is_pii(self, value: str) -> bool:
        """Detect if a string value looks like PII"""
        # Simple heuristics for PII detection
        pii_patterns = [
            r'\b\d{3}-\d{2}-\d{4}\b',  # SSN
            r'\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b',  # Credit card
            r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b',  # Email
        ]
        for pattern in pii_patterns:
            if re.search(pattern, value):
                return True
        return False

    def _hash_pii(self, value: str) -> str:
        """Hash PII value for storage"""
        return f"[HASH:{hashlib.sha256(value.encode()).hexdigest()[:16]}]"

    def capture_event(self, event_type: str, data: Dict[str, Any], trace_id: Optional[str] = None):
        """Capture a security event"""
        event = {
            "type": event_type,
            "data": self.redact_data(data),
            "timestamp": self._current_timestamp(),
            "service": self.config.service_name,
            "env": self.config.env,
            "trace_id": trace_id,
        }

        if self.config.offline_mode:
            self._cache_event(event)
        else:
            self.event_queue.append(event)
            if len(self.event_queue) >= self.config.telemetry_batch_size:
                self._flush_events()

    def _current_timestamp(self) -> int:
        """Get current timestamp in milliseconds"""
        import time
        return int(time.time() * 1000)

    def _cache_event(self, event: Dict[str, Any]):
        """Cache event for offline mode"""
        import os
        cache_dir = self.config.cache_dir or os.path.join(os.getcwd(), ".security-platform-cache")
        os.makedirs(cache_dir, exist_ok=True)
        cache_file = os.path.join(cache_dir, "events.jsonl")
        with open(cache_file, "a") as f:
            f.write(json.dumps(event) + "\n")

    def _flush_events(self):
        """Flush queued events to control plane"""
        if not self.event_queue:
            return

        import requests
        batch = self.event_queue[:]
        self.event_queue = []

        try:
            response = requests.post(
                f"{self.config.control_plane_url}/v1/events",
                json={"events": batch},
                headers={
                    "Authorization": f"Bearer {self.config.auth_key}",
                    "Content-Type": "application/json",
                },
                timeout=5,
            )
            response.raise_for_status()
        except Exception as e:
            # Re-queue events on failure
            self.event_queue = batch + self.event_queue
            print(f"Failed to send events: {e}")

    def get_tracer(self, name: str):
        """Get an OpenTelemetry tracer"""
        return trace.get_tracer(name)

    def shutdown(self):
        """Shutdown the agent"""
        self._flush_events()
        if self.tracer_provider:
            self.tracer_provider.shutdown()


# Global agent instance
_agent: Optional[SecurityAgent] = None


def init_agent(config: Optional[AgentConfig] = None) -> SecurityAgent:
    """Initialize the global agent instance"""
    global _agent
    if config is None:
        config = AgentConfig.from_env()
    _agent = SecurityAgent(config)
    return _agent


def get_agent() -> Optional[SecurityAgent]:
    """Get the global agent instance"""
    return _agent
