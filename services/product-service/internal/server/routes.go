package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerRoutes() {
	// Health & readiness probes — framework layer, lives here.
	s.echo.GET("/health", s.handleHealth)
	s.echo.GET("/ready", s.handleReady)

	// Feature module routes are registered from main.go via srv.Group()
	// so that each module owns its own route registration.
}

func (s *Server) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": s.cfg.App.Name,
	})
}

func (s *Server) handleReady(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}
