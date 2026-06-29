package domain

import "time"

// Document represents a gated document with a price.
type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	PriceUSD  int32     `json:"price_usd"` // in cents
	CreatedAt time.Time `json:"created_at"`
}
