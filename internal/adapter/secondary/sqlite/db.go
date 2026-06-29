package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// DB represents a connection to the SQLite database.
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection.
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	return &DB{db}, nil
}

// RunMigrations runs Goose database migrations with the specified command (e.g. "up", "down", "status").
func RunMigrations(db *sql.DB, migrationsFS fs.FS, dir string, command string) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	var err error
	switch command {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "status":
		err = goose.Status(db, dir)
	default:
		return fmt.Errorf("unsupported migration command: %s", command)
	}

	if err != nil {
		return fmt.Errorf("goose %s failed on directory %s: %w", command, dir, err)
	}

	return nil
}

// Ping checks if the database is reachable.
func (d *DB) PingContext(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}
