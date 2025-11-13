"""FastAPI/Starlette middleware for security agent"""

from typing import Callable
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.requests import Request
from starlette.responses import Response
from opentelemetry import trace

from .agent import get_agent


class SecurityMiddleware(BaseHTTPMiddleware):
    """Middleware for capturing requests and security events"""

    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        agent = get_agent()
        if not agent:
            return await call_next(request)

        tracer = agent.get_tracer("security-middleware")
        span = tracer.start_span(f"{request.method} {request.url.path}")

        try:
            # Capture request metadata
            request_data = {
                "method": request.method,
                "path": request.url.path,
                "query_params": dict(request.query_params),
                "headers": dict(request.headers),
                "client": request.client.host if request.client else None,
            }

            # Redact sensitive headers
            request_data["headers"] = agent.redact_data(request_data["headers"])

            # Process request
            response = await call_next(request)

            # Capture response metadata
            response_data = {
                "status_code": response.status_code,
                "headers": dict(response.headers),
            }

            # Capture security event
            event_data = {
                "request": request_data,
                "response": response_data,
            }

            # Check for potential security issues
            if response.status_code >= 400:
                agent.capture_event("http_error", event_data, span.get_span_context().trace_id)

            # Detect suspicious patterns
            if self._is_suspicious(request):
                agent.capture_event("suspicious_request", event_data, span.get_span_context().trace_id)

            span.set_attribute("http.status_code", response.status_code)
            span.set_attribute("http.method", request.method)
            span.set_attribute("http.url", str(request.url))

            return response

        except Exception as e:
            agent.capture_event("request_exception", {
                "error": str(e),
                "path": request.url.path,
            }, span.get_span_context().trace_id if span else None)
            span.record_exception(e)
            raise
        finally:
            span.end()

    def _is_suspicious(self, request: Request) -> bool:
        """Detect suspicious request patterns"""
        path = request.url.path.lower()
        suspicious_patterns = [
            "/admin",
            "/.env",
            "/wp-admin",
            "/phpmyadmin",
            "..",
            "cmd=",
            "exec=",
            "eval(",
        ]

        for pattern in suspicious_patterns:
            if pattern in path:
                return True

        # Check for SQL injection patterns in query params
        for param, value in request.query_params.items():
            if any(keyword in value.lower() for keyword in ["union", "select", "drop", "delete"]):
                return True

        return False
