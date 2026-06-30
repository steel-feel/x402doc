package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path"

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
// It executes migrations from the "ddl" subfolder first, followed by "dml".
func RunMigrations(db *sql.DB, migrationsFS fs.FS, dir string, command string) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	ddlDir := path.Join(dir, "ddl")
	dmlDir := path.Join(dir, "dml")

	switch command {
	case "up":
		// 1. Run DDL migrations
		goose.SetTableName("goose_db_version_ddl")
		if err := goose.Up(db, ddlDir); err != nil {
			return fmt.Errorf("failed to run goose DDL migrations: %w", err)
		}

		// 2. Run DML migrations
		goose.SetTableName("goose_db_version_dml")
		if err := goose.Up(db, dmlDir); err != nil {
			return fmt.Errorf("failed to run goose DML migrations: %w", err)
		}

	case "down":
		// 1. Rollback DML migrations first (since DML data depends on DDL schema)
		goose.SetTableName("goose_db_version_dml")
		if err := goose.Down(db, dmlDir); err != nil {
			return fmt.Errorf("failed to rollback goose DML migrations: %w", err)
		}

		// 2. Rollback DDL migrations
		goose.SetTableName("goose_db_version_ddl")
		if err := goose.Down(db, ddlDir); err != nil {
			return fmt.Errorf("failed to rollback goose DDL migrations: %w", err)
		}

	case "status":
		// Print DDL status
		fmt.Println("--- DDL Migration Status ---")
		goose.SetTableName("goose_db_version_ddl")
		if err := goose.Status(db, ddlDir); err != nil {
			return fmt.Errorf("failed to get DDL status: %w", err)
		}

		// Print DML status
		fmt.Println("\n--- DML Migration Status ---")
		goose.SetTableName("goose_db_version_dml")
		if err := goose.Status(db, dmlDir); err != nil {
			return fmt.Errorf("failed to get DML status: %w", err)
		}

	default:
		return fmt.Errorf("unsupported migration command: %s", command)
	}

	return nil
}

// Ping checks if the database is reachable.
func (d *DB) PingContext(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}
