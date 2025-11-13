"""
Security Platform Python Agent
"""

from .agent import SecurityAgent, SecurityMiddleware, create_agent
from .config import AgentConfig

__version__ = "0.1.0"
__all__ = ["SecurityAgent", "SecurityMiddleware", "create_agent", "AgentConfig"]
