package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

// ObservabilityMiddleware sets up structured logging and can be extended with metrics/tracing.
func ObservabilityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)
			slog.Info("request", "method", c.Request().Method, "path", c.Request().URL.Path, "duration", time.Since(start))
			return err
		}
	}
}
