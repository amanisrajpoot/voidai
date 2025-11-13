"""Middleware for FastAPI and Django."""

import time
from typing import Callable, Optional
from .agent import Agent

try:
    from fastapi import Request, Response
    from starlette.middleware.base import BaseHTTPMiddleware
    FASTAPI_AVAILABLE = True
except ImportError:
    FASTAPI_AVAILABLE = False


def create_middleware(
    agent: Agent, mode: str = "observe", capture_response_body: bool = False
) -> Optional[Callable]:
    """Create middleware for FastAPI or Django."""
    if FASTAPI_AVAILABLE:
        return FastAPIMiddleware(agent, mode, capture_response_body)
    # Django middleware would go here
    return None


if FASTAPI_AVAILABLE:

    class FastAPIMiddleware(BaseHTTPMiddleware):
        """FastAPI middleware for security observability."""

        def __init__(
            self,
            agent: Agent,
            mode: str = "observe",
            capture_response_body: bool = False,
        ):
            super().__init__(None)
            self.agent = agent
            self.mode = mode
            self.capture_response_body = capture_response_body

        async def dispatch(self, request: Request, call_next):
            """Process request and response."""
            start_time = time.time()

            # Capture request metadata
            request_metadata = {
                "method": request.method,
                "path": request.url.path,
                "headers": self.agent.redact(dict(request.headers)),
                "query": self.agent.redact(dict(request.query_params)),
                "ip": request.client.host if request.client else None,
            }

            # Check policy
            policy_result = self.agent.check_policy(
                "http_request",
                {
                    "method": request.method,
                    "path": request.url.path,
                    "ip": request.client.host if request.client else None,
                },
            )

            if not policy_result["allowed"] and self.mode == "block":
                from fastapi.responses import JSONResponse
                return JSONResponse(
                    status_code=403,
                    content={
                        "error": "Request blocked by security policy",
                        "reason": policy_result.get("reason"),
                    },
                )

            # Process request
            response = await call_next(request)

            # Capture response metadata
            duration = time.time() - start_time
            response_metadata = {
                "status_code": response.status_code,
                "headers": self.agent.redact(dict(response.headers)),
                "duration": duration,
            }

            # Send event (async, non-blocking)
            # In real implementation, send to control plane

            return response
