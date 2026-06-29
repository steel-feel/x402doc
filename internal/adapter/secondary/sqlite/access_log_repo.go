package sqlite

import (
	"context"
	"fmt"

	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

type accessLogRepo struct {
	db *DB
}

// NewAccessLogRepository creates a new SQLite-backed access log repository.
func NewAccessLogRepository(db *DB) port.AccessLogRepository {
	return &accessLogRepo{db: db}
}

func (r *accessLogRepo) Save(ctx context.Context, log *domain.AccessLog) error {
	_, err := r.db.ExecContext(ctx, 
		"INSERT INTO access_logs (document_id, blockchain_address, ip_address, user_agent, payment_tx_hash) VALUES (?, ?, ?, ?, ?)",
		log.DocumentID, log.BlockchainAddress, log.IPAddress, log.UserAgent, log.PaymentTxHash)
	if err != nil {
		return fmt.Errorf("failed to insert access log: %w", err)
	}
	return nil
}
