-- documents table
CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    price_usd INTEGER NOT NULL,  -- price in cents (e.g., 100 = $1.00)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- access_logs table
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

-- seed data
INSERT OR IGNORE INTO documents (id, title, content, price_usd) VALUES
    ('doc-001', 'Getting Started with Tempo', 'Full guide to Tempo blockchain...', 100),
    ('doc-002', 'TIP-20 Token Standard', 'Deep dive into TIP-20...', 250);
