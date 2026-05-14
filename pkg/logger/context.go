package logger

import (
	"context"

	"go.uber.org/zap"
)

// contextKey is an unexported type to prevent key collisions with other
// packages that also store values in context.
type contextKey struct{}

// WithContext returns a new context carrying the given logger.
// Call this once per request in middleware, with a logger already enriched
// with request_id, method, and URI — so every downstream log is pre-tagged.
func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

// FromContext retrieves the request-scoped logger from ctx.
// Falls back to the global logger (set via zap.ReplaceGlobals in main.go)
// so it is always safe to call — even outside an HTTP request context.
func FromContext(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return log
	}
	return zap.L() // global fallback
}
