"""Type definitions for the agent."""

from typing import Dict, Any, List, Optional, Literal

AgentConfig = Dict[str, Any]


class RedactionRule:
    """Rule for redacting sensitive data."""

    def __init__(
        self,
        pattern: str,
        replacement: Optional[str] = None,
        field: Optional[str] = None,
    ):
        self.pattern = pattern
        self.replacement = replacement or "[REDACTED]"
        self.field = field  # "*" for all fields, or specific field name


class PolicyRule:
    """Policy rule for access control."""

    def __init__(
        self,
        id: str,
        action: str,
        condition: str,
        effect: Literal["allow", "deny"],
    ):
        self.id = id
        self.action = action
        self.condition = condition
        self.effect = effect
