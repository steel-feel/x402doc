package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/api/generated"
	"github.com/steel-feel/prac/internal/port"
)

type DocumentHandler struct {
	service port.DocumentService
}

func NewDocumentHandler(s port.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: s}
}

func (h *DocumentHandler) GetDoc(c *echo.Context) error {
	id := c.Param("id")
	doc, err := h.service.GetDocument(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get document"})
	}
	if doc == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "document not found"})
	}

	resp := generated.Document{
		Id:      doc.ID,
		Title:   doc.Title,
		Content: doc.Content,
		Price:   doc.PriceUSD,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *DocumentHandler) CreateDoc(c *echo.Context) error {
	var req generated.CreateDocumentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Title == "" || req.Content == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title, content, and positive price are required"})
	}

	doc, err := h.service.CreateDocument(c.Request().Context(), req.Title, req.Price, req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create document"})
	}

	resp := generated.Document{
		Id:      doc.ID,
		Title:   doc.Title,
		Content: doc.Content,
		Price:   doc.PriceUSD,
	}

	return c.JSON(http.StatusCreated, resp)
}
