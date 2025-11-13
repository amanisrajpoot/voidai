"""
Security Platform Python Agent

Provides RASP (Runtime Application Self-Protection) and observability
for Python applications using FastAPI, Django, Flask, etc.
"""

from .agent import Agent
from .types import AgentConfig, RedactionRule, PolicyRule
from .middleware import create_middleware

__all__ = ["Agent", "AgentConfig", "RedactionRule", "PolicyRule", "create_middleware"]
__version__ = "1.0.0"
