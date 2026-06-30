package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/api/generated"
	"github.com/steel-feel/prac/internal/port"
)

type PriceHandler struct {
	service port.PriceService
}

// NewPriceHandler creates a new REST HTTP handler for market price queries.
func NewPriceHandler(s port.PriceService) *PriceHandler {
	return &PriceHandler{service: s}
}

// GetEthereumPrice handles the GET /api/v1/price/eth REST call.
func (h *PriceHandler) GetEthereumPrice(c *echo.Context) error {
	price, err := h.service.GetEthereumPrice(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch ethereum price"})
	}

	resp := generated.EthereumPriceResponse{
		Price: price,
	}

	return c.JSON(http.StatusOK, resp)
}
