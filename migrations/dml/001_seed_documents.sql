-- +goose Up
-- +goose StatementBegin
INSERT OR IGNORE INTO documents (id, title, content, price_usd) VALUES
    ('doc-001', 'Getting Started with Tempo', 'Full guide to Tempo blockchain...', 100),
    ('doc-002', 'TIP-20 Token Standard', 'Deep dive into TIP-20...', 250);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM documents WHERE id IN ('doc-001', 'doc-002');
-- +goose StatementEnd
