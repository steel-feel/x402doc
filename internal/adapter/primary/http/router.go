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
	Auth     *AuthHandler
	Price    *PriceHandler
}

// RegisterRoutes sets up all the routes for the Echo application.
func RegisterRoutes(e *echo.Echo, h *Handlers, fac *facilitator.Facilitator, docSvc port.DocumentService, accessRepo port.AccessLogRepository, jwtSecret string) {
	api := e.Group("/api/v1")

	api.GET("/health", h.Health.Check)
	api.POST("/login", h.Auth.Login)

	// Apply x402 middleware to document endpoints
	x402Mw := middleware.X402Middleware(docSvc, fac)
	accessLogMw := middleware.AccessLogMiddleware(accessRepo)
	
	api.GET("/documents/:id", h.Document.GetDoc, x402Mw, accessLogMw)

	// Apply JWT authentication middleware to document creation route
	api.POST("/documents", h.Document.CreateDoc, middleware.JWTMiddleware(jwtSecret))

	// Ethereum price query endpoint
	api.GET("/price/eth", h.Price.GetEthereumPrice)
}
