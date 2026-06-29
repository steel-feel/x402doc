package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

type documentRepo struct {
	db *DB
}

// NewDocumentRepository creates a new SQLite-backed document repository.
func NewDocumentRepository(db *DB) port.DocumentRepository {
	return &documentRepo{db: db}
}

func (r *documentRepo) FindByID(ctx context.Context, id string) (*domain.Document, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, title, content, price_usd, created_at FROM documents WHERE id = ?", id)

	var doc domain.Document
	err := row.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.PriceUSD, &doc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Return nil, nil for not found (or return a custom error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find document by id: %w", err)
	}

	return &doc, nil
}

func (r *documentRepo) List(ctx context.Context) ([]*domain.Document, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, content, price_usd, created_at FROM documents")
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var docs []*domain.Document
	for rows.Next() {
		var doc domain.Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.PriceUSD, &doc.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		docs = append(docs, &doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return docs, nil
}
