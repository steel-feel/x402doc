package port

import (
	"context"

	"github.com/steel-feel/prac/internal/domain"
)

// DocumentRepository is the secondary port for document storage.
type DocumentRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Document, error)
	List(ctx context.Context) ([]*domain.Document, error)
}

// AccessLogRepository is the secondary port for storing access logs.
type AccessLogRepository interface {
	Save(ctx context.Context, log *domain.AccessLog) error
}
