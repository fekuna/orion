package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/fekuna/orion-v2/services/product-service/internal/config"
)

// customValidator bridges go-playground/validator with Echo's Validator interface.
type customValidator struct {
	v *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

const (
	shutdownTimeout = 10 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 30 * time.Second
)

// Server wraps the Echo instance and service dependencies.
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	log  *zap.Logger
}

// New creates and configures a new Echo server.
// The provided logger is used for both application events and HTTP access logs.
func New(cfg *config.Config, log *zap.Logger) *Server {
	e := echo.New()

	// --- global settings ---
	e.HideBanner = true
	e.HidePort = true
	e.Debug = cfg.App.Debug

	// Replace Echo's built-in logger with a no-op so Zap is the only
	// logging path (avoids duplicate output).
	e.Logger.SetOutput(io.Discard)

	// Register the request validator.
	e.Validator = &customValidator{v: validator.New()}

	// Register the centralized HTTP error handler.
	// All handlers return echo.NewHTTPError — this formats them into
	// the standard response envelope and logs 5xx errors with request context.
	e.HTTPErrorHandler = httpErrorHandler(log)

	// --- middleware stack ---
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error("panic recovered",
				zap.Error(err),
				zap.ByteString("stack", stack),
			)
			return nil
		},
	}))
	e.Use(middleware.RequestID())
	e.Use(zapAccessLogger(log))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.HTTP.AllowedOrigins},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: writeTimeout,
	}))

	s := &Server{echo: e, cfg: cfg, log: log}
	s.registerRoutes()

	return s
}

// Start begins listening and handles graceful shutdown on OS signals.
func (s *Server) Start() error {
	// Configure underlying http.Server timeouts.
	s.echo.Server.ReadTimeout = readTimeout
	s.echo.Server.WriteTimeout = writeTimeout

	addr := s.cfg.HTTP.Addr()

	// Start in a goroutine so we can listen for shutdown signals.
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("server starting",
			zap.String("address", addr),
			zap.String("env", s.cfg.App.Env),
		)
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Block until a signal or a startup error arrives.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		s.log.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}

	s.log.Info("server stopped gracefully")
	return nil
}

// Group creates a new Echo router group with an optional prefix and middleware.
// Use this in main.go to register module routes after server creation.
func (s *Server) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return s.echo.Group(prefix, m...)
}

// zapAccessLogger returns an Echo middleware that writes structured HTTP
// access log entries using Zap.
func zapAccessLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			req := c.Request()
			res := c.Response()

			fields := []zap.Field{
				zap.String("id", res.Header().Get(echo.HeaderXRequestID)),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("remote_ip", c.RealIP()),
				zap.Int64("bytes_in", req.ContentLength),
				zap.Int64("bytes_out", res.Size),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				log.Error("request", fields...)
			} else {
				log.Info("request", fields...)
			}

			return err
		}
	}
}
