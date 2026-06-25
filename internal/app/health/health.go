package health

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	// domain services here
}

func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes ties this controller to the Echo router
func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.GET("/health", h.Check)
}

func (h *Handler) Check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status": "fully operational",
	})
}
