package domain

import "time"

// AccessLog represents a record of a document being accessed.
type AccessLog struct {
	ID                int64     `json:"id"`
	DocumentID        string    `json:"document_id"`
	BlockchainAddress string    `json:"blockchain_address"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	PaymentTxHash     string    `json:"payment_tx_hash"`
	AccessedAt        time.Time `json:"accessed_at"`
}
