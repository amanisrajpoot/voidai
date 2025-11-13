from typing import List, Literal, Optional
from dataclasses import dataclass, field


@dataclass
class AgentConfig:
    """Configuration for Security Platform Agent"""
    control_plane_url: str
    auth_key: str
    service_name: Optional[str] = None
    environment: Optional[str] = None
    redaction_rules: List[str] = field(default_factory=lambda: [
        "password", "card", "ssn", "pin", "auth.*"
    ])
    local_policy: Literal["observe", "block"] = "observe"
    telemetry_batch_size: int = 50
    otlp_endpoint: Optional[str] = None

    def __post_init__(self):
        if self.service_name is None:
            self.service_name = "python-service"
        if self.environment is None:
            self.environment = "production"
        if self.otlp_endpoint is None:
            self.otlp_endpoint = f"{self.control_plane_url}/v1/traces"
