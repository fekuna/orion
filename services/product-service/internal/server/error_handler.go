package server

import (
	"errors"
	"net/http"

	"github.com/fekuna/orion-v2/pkg/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// httpErrorHandler is Echo's centralized error handler.
// All handlers return echo.NewHTTPError — this function formats
// them into the standard response envelope.
//
// Benefits:
//   - Error formatting in one place (DRY)
//   - Request context (ID, method, URI) always present in logs
//   - Handlers stay thin — just return the right error type
func httpErrorHandler(log *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var he *echo.HTTPError
		if !errors.As(err, &he) {
			he = &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		// Log 5xx errors with full request context for correlation.
		if he.Code >= http.StatusInternalServerError {
			log.Error("server error",
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
				zap.Int("status", he.Code),
				zap.Error(err),
			)
		}

		var body any
		switch he.Code {
		case http.StatusBadRequest:
			body = response.ErrBadRequest(stringMessage(he.Message))
		case http.StatusUnauthorized:
			body = response.ErrUnauthorized()
		case http.StatusForbidden:
			body = response.ErrForbidden()
		case http.StatusNotFound:
			body = response.ErrNotFound(stringMessage(he.Message))
		case http.StatusConflict:
			body = response.ErrConflict(stringMessage(he.Message))
		case http.StatusUnprocessableEntity:
			body = response.ErrValidation(he.Message)
		default:
			body = response.ErrInternal()
		}

		_ = c.JSON(he.Code, body)
	}
}

// stringMessage safely converts an echo.HTTPError message to a string.
func stringMessage(msg any) string {
	if s, ok := msg.(string); ok {
		return s
	}
	return "unknown error"
}
