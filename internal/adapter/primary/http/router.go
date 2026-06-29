package http

import (
	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/internal/adapter/primary/http/middleware"
	"github.com/steel-feel/prac/internal/adapter/secondary/facilitator"
	"github.com/steel-feel/prac/internal/port"
)

// Handlers holds all HTTP handlers for the application.
type Handlers struct {
	Health   *HealthHandler
	Document *DocumentHandler
}

// RegisterRoutes sets up all the routes for the Echo application.
func RegisterRoutes(e *echo.Echo, h *Handlers, fac *facilitator.Facilitator, docSvc port.DocumentService, accessRepo port.AccessLogRepository) {
	api := e.Group("/api/v1")

	api.GET("/health", h.Health.Check)

	// Apply x402 middleware to document endpoints
	x402Mw := middleware.X402Middleware(docSvc, fac)
	accessLogMw := middleware.AccessLogMiddleware(accessRepo)
	
	api.GET("/documents/:id", h.Document.GetDoc, x402Mw, accessLogMw)
}
