package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/api/generated"
	"github.com/steel-feel/prac/internal/port"
)

type HealthHandler struct {
	service port.HealthService
}

func NewHealthHandler(s port.HealthService) *HealthHandler {
	return &HealthHandler{service: s}
}

func (h *HealthHandler) Check(c *echo.Context) error {
	status, err := h.service.Check(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check health"})
	}

	resp := generated.HealthStatus{
		Status:   status.Status,
		DbStatus: status.DBStatus,
		Uptime:   status.Uptime,
	}

	return c.JSON(http.StatusOK, resp)
}
