package httputil

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ParseUUIDParam extracts and parses a UUID path parameter by name.
// Returns (uuid.Nil, echo.HTTPError 400) if the value cannot be parsed —
// the caller can return this error directly to Echo.
func ParseUUIDParam(c echo.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid "+name)
	}
	return id, nil
}
