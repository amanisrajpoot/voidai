"""PII redaction implementation."""

import re
from typing import Any, List, Optional
from .types import RedactionRule

DEFAULT_BLOCKLIST = [
    "password",
    "passwd",
    "pwd",
    "secret",
    "token",
    "api_key",
    "apikey",
    "auth",
    "credit_card",
    "card_number",
    "cvv",
    "ssn",
    "social_security",
    "pin",
    "passcode",
]


class Redactor:
    """Redacts PII from data structures."""

    def __init__(self, rules: Optional[List[RedactionRule]] = None):
        """Initialize redactor with rules."""
        self.rules = rules or []
        self.patterns = [
            re.compile(rf"\b{field}\b", re.IGNORECASE) for field in DEFAULT_BLOCKLIST
        ]

    def redact(self, obj: Any) -> Any:
        """Recursively redact PII from object."""
        if obj is None:
            return obj

        if isinstance(obj, str):
            return self._redact_string(obj)

        if isinstance(obj, list):
            return [self.redact(item) for item in obj]

        if isinstance(obj, dict):
            redacted = {}
            for key, value in obj.items():
                redacted_key = (
                    "[REDACTED]" if self._should_redact_key(key) else key
                )
                redacted[redacted_key] = self.redact(value)
            return redacted

        return obj

    def _should_redact_key(self, key: str) -> bool:
        """Check if a key should be redacted."""
        lower_key = key.lower()

        # Check default blocklist
        for pattern in self.patterns:
            if pattern.search(lower_key):
                return True

        # Check custom rules
        for rule in self.rules:
            if rule.field == "*" or rule.field == key:
                try:
                    regex = re.compile(rule.pattern, re.IGNORECASE)
                    if regex.search(key) or regex.search(lower_key):
                        return True
                except re.error:
                    # Invalid regex, skip
                    pass

        return False

    def _redact_string(self, s: str) -> str:
        """Redact PII from string."""
        for rule in self.rules:
            if not rule.field or rule.field == "*":
                try:
                    regex = re.compile(rule.pattern, re.IGNORECASE)
                    if regex.search(s):
                        return regex.sub(rule.replacement, s)
                except re.error:
                    # Invalid regex, skip
                    pass
        return s
