package service

import (
	"context"

	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

type documentService struct {
	repo port.DocumentRepository
}

// NewDocumentService creates a new document service.
func NewDocumentService(repo port.DocumentRepository) port.DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*domain.Document, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *documentService) ListDocuments(ctx context.Context) ([]*domain.Document, error) {
	return s.repo.List(ctx)
}
