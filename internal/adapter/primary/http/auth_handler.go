package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/api/generated"
	"github.com/steel-feel/prac/internal/port"
)

type AuthHandler struct {
	service port.UserService
}

func NewAuthHandler(s port.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req generated.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	}

	token, err := h.service.Authenticate(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	resp := generated.LoginResponse{
		Token: token,
	}

	return c.JSON(http.StatusOK, resp)
}
