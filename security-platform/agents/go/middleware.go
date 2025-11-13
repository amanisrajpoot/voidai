package agent

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware creates HTTP middleware for automatic instrumentation.
func (a *Agent) HTTPMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			span := trace.SpanFromContext(ctx)

			// Add request attributes
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
			)

			// Check policy if in block mode
			if a.config.Policy != nil && a.config.Policy.Mode == "block" {
				result := a.CheckPolicy("http.request", map[string]interface{}{
					"method": r.Method,
					"path":   r.URL.Path,
					"ip":     r.RemoteAddr,
				})

				if !result.Allowed {
					span.SetAttributes(attribute.Bool("policy.blocked", true))
					span.SetStatus(codes.Error, result.Reason)
					http.Error(w, "Request blocked by policy", http.StatusForbidden)
					return
				}
			}

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		}),
		"http.request",
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
	)
}
