package agent

import (
	"context"
	"net/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Middleware provides HTTP middleware for the agent
type Middleware struct {
	agent       *Agent
	redaction   *Redaction
	policyEngine *PolicyEngine
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(agent *Agent) *Middleware {
	return &Middleware{
		agent:       agent,
		redaction:   NewRedaction(nil),
		policyEngine: NewPolicyEngine(agent.GetConfig()),
	}
}

// Handler wraps an HTTP handler with security platform instrumentation
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("security-platform-agent")

		ctx, span := tracer.Start(ctx, "http.request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
			))

		defer span.End()

		// Capture request data
		requestData := map[string]interface{}{
			"method": r.Method,
			"uri":    r.URL.Path,
			"query":  r.URL.RawQuery,
		}

		headers := make(map[string]interface{})
		for k, v := range r.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		requestData["headers"] = m.redaction.RedactMap(headers)

		// Check policy
		policyResult := m.policyEngine.CheckPolicy("http_request", requestData)

		if m.agent.GetConfig().LocalPolicy == "block" && !policyResult.IsAllowed {
			http.Error(w, "Request blocked by security policy", http.StatusForbidden)
			span.SetAttributes(
				attribute.Bool("security.blocked", true),
				attribute.String("security.reason", policyResult.Reason),
			)
			return
		}

		// Create response writer wrapper
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r.WithContext(ctx))

		span.SetAttributes(
			attribute.Int("http.status_code", wrapped.statusCode),
			attribute.Bool("security.policy.allowed", policyResult.IsAllowed),
		)

		if !policyResult.IsAllowed {
			span.SetAttributes(attribute.String("security.policy.reason", policyResult.Reason))
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
