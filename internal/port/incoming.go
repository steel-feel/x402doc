package port

import (
	"context"

	"github.com/steel-feel/prac/internal/domain"
)

// DocumentService is the primary port for document operations.
type DocumentService interface {
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	ListDocuments(ctx context.Context) ([]*domain.Document, error)
	CreateDocument(ctx context.Context, title string, price int32, content string) (*domain.Document, error)
}

// HealthService is the primary port for checking service health.
type HealthService interface {
	Check(ctx context.Context) (*domain.HealthStatus, error)
}

// UserService is the primary port for user authentication.
type UserService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

// PriceService is the primary port for market price operations.
type PriceService interface {
	GetEthereumPrice(ctx context.Context) (float64, error)
}
