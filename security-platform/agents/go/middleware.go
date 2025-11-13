package agent

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// Middleware creates HTTP middleware for the agent.
func (a *Agent) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check policy
		context := map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"remote_addr": r.RemoteAddr,
		}

		result := a.CheckPolicy("http_request", context)
		if !result.Allowed {
			a.logger.Warn("Request blocked by policy",
				zap.String("reason", result.Reason),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			http.Error(w, "Request blocked by security policy", http.StatusForbidden)
			return
		}

		// Continue with request
		next.ServeHTTP(w, r)
	})
}

// HTTPClient returns an HTTP client with OpenTelemetry instrumentation.
func (a *Agent) HTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}
