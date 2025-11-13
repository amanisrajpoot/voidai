"""Policy engine for access control."""

from typing import Dict, Any, List, Optional
from .types import PolicyRule


class PolicyEngine:
    """Evaluates policy rules."""

    def __init__(self, config: Dict[str, Any]):
        """Initialize policy engine."""
        self.mode = config.get("mode", "observe")
        self.rules: List[PolicyRule] = config.get("rules", [])

    def check(
        self, action: str, context: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Check if action is allowed."""
        # In observe mode, always allow but log
        if self.mode == "observe":
            return {"allowed": True}

        # Check rules
        for rule in self.rules:
            if rule.action == action or rule.action == "*":
                if self._evaluate_condition(rule.condition, context):
                    return {
                        "allowed": rule.effect == "allow",
                        "reason": rule.id,
                    }

        # Default allow if no rules match
        return {"allowed": True}

    def _evaluate_condition(
        self, condition: str, context: Dict[str, Any]
    ) -> bool:
        """Evaluate condition expression."""
        try:
            # Simple key-value matching
            # In real implementation, use a proper expression evaluator
            parts = condition.split("==")
            if len(parts) == 2:
                key = parts[0].strip()
                value = parts[1].strip().strip("'\"")
                return context.get(key) == value
            return False
        except Exception:
            return False
