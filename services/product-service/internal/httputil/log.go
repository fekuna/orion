// Package httputil provides HTTP-related utilities shared across all
// feature modules within this service.
//
// This package is internal to the service — use pkg/ for utilities
// that need to be shared across multiple services.
package httputil

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LogError logs a handler-level error with full request context for log correlation.
// Always call this before returning a 5xx echo.HTTPError from a handler so the
// error entry can be matched to its corresponding access log entry via request_id.
func LogError(log *zap.Logger, c echo.Context, op string, err error, extra ...zap.Field) {
	fields := []zap.Field{
		zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		zap.String("method", c.Request().Method),
		zap.String("uri", c.Request().RequestURI),
		zap.String("op", op),
		zap.Error(err),
	}
	log.Error("handler error", append(fields, extra...)...)
}
