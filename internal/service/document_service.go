package service

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (s *documentService) CreateDocument(ctx context.Context, title string, price int32, content string) (*domain.Document, error) {
	doc := &domain.Document{
		ID:        uuid.NewString(),
		Title:     title,
		PriceUSD:  price,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}
