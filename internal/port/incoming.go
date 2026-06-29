package port

import (
	"context"

	"github.com/steel-feel/prac/internal/domain"
)

// DocumentService is the primary port for document operations.
type DocumentService interface {
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	ListDocuments(ctx context.Context) ([]*domain.Document, error)
}

// HealthService is the primary port for checking service health.
type HealthService interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}
