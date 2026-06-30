-- +goose Up
-- +goose StatementBegin
-- Password is "password" (bcrypt hash)
INSERT OR IGNORE INTO users (id, username, password_hash) VALUES
    ('user-001', 'admin', '$2a$10$f.xS9f4m.w8CMw7oa6tPcOUC7fgaQrX7Bo12AXqBoVl1Yzni5819S');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 'user-001';
-- +goose StatementEnd
