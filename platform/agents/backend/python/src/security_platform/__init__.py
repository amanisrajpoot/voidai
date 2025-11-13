"""
Security Platform Agent for Python
Provides RASP-lite capabilities, telemetry collection, and security monitoring
"""

from .agent import SecurityAgent, init_agent
from .middleware import SecurityMiddleware
from .config import AgentConfig

__version__ = "0.1.0"
__all__ = ["SecurityAgent", "init_agent", "SecurityMiddleware", "AgentConfig"]
