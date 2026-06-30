package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

type userRepo struct {
	db *DB
}

// NewUserRepository creates a new SQLite-backed user repository.
func NewUserRepository(db *DB) port.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, username, password_hash, created_at FROM users WHERE username = ?", username)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &u, nil
}
