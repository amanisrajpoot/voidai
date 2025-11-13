"""Configuration management for the security agent"""

import os
from pathlib import Path
from typing import List, Optional
from pydantic_settings import BaseSettings
from pydantic import Field
import yaml


class AgentConfig(BaseSettings):
    """Agent configuration"""

    # Control plane settings
    control_plane_url: str = Field(
        default="https://api.example.com",
        description="Control plane API endpoint"
    )
    auth_key: str = Field(
        default="",
        description="Authentication key for control plane"
    )

    # Environment
    env: str = Field(
        default="development",
        description="Deployment environment"
    )
    service_name: str = Field(
        default="python-service",
        description="Service identifier"
    )

    # Redaction rules
    redaction_rules: List[str] = Field(
        default_factory=lambda: [
            "password",
            "card",
            "ssn",
            "pin",
            "auth.*",
            "token",
            "secret",
            "api[_-]?key",
        ],
        description="Patterns to redact from telemetry"
    )

    # Local policy
    local_policy_mode: str = Field(
        default="observe",
        description="Policy mode: observe or block"
    )

    # Telemetry settings
    telemetry_batch_size: int = Field(
        default=100,
        description="Number of events per batch"
    )
    otlp_endpoint: str = Field(
        default="http://localhost:4318",
        description="OTLP exporter endpoint"
    )

    # Offline mode
    offline_mode: bool = Field(
        default=False,
        description="Enable offline mode (cache telemetry locally)"
    )
    cache_dir: Optional[str] = Field(
        default=None,
        description="Directory for offline cache"
    )

    @classmethod
    def from_yaml(cls, path: str) -> "AgentConfig":
        """Load configuration from YAML file"""
        with open(path, "r") as f:
            data = yaml.safe_load(f)
        return cls(**data)

    @classmethod
    def from_env(cls) -> "AgentConfig":
        """Load configuration from environment variables"""
        return cls(
            control_plane_url=os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL", "https://api.example.com"),
            auth_key=os.getenv("SECURITY_PLATFORM_AUTH_KEY", ""),
            env=os.getenv("SECURITY_PLATFORM_ENV", "development"),
            service_name=os.getenv("SECURITY_PLATFORM_SERVICE_NAME", "python-service"),
            otlp_endpoint=os.getenv("SECURITY_PLATFORM_OTLP_ENDPOINT", "http://localhost:4318"),
            offline_mode=os.getenv("SECURITY_PLATFORM_OFFLINE_MODE", "false").lower() == "true",
        )

    class Config:
        env_prefix = "SECURITY_PLATFORM_"
        case_sensitive = False
