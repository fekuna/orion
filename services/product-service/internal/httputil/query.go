package httputil

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// ParseIntQuery reads an integer query parameter, returning defaultVal
// when the parameter is absent, non-numeric, or <= 0.
func ParseIntQuery(c echo.Context, key string, defaultVal int) int {
	v := c.QueryParam(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
