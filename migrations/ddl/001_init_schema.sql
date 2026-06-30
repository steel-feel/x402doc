-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    price_usd INTEGER NOT NULL,  -- price in cents (e.g., 100 = $1.00)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS access_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id TEXT NOT NULL,
    blockchain_address TEXT,
    ip_address TEXT,
    user_agent TEXT,
    payment_tx_hash TEXT,
    accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_logs;
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
