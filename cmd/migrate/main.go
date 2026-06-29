package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/steel-feel/prac/internal/adapter/secondary/sqlite"
)

func main() {
	dbPath := flag.String("db-path", "data.db", "path to sqlite database")
	flag.Parse()

	// Default command is "up"
	cmd := "up"
	args := flag.Args()
	if len(args) > 0 {
		cmd = args[0]
	}

	// 1. Initialize DB connection
	db, err := sqlite.NewDB(*dbPath)
	if err != nil {
		fmt.Printf("Error: failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 2. Run migrations
	migrationsFS := os.DirFS(".")
	fmt.Printf("Executing migration command '%s' on database '%s'...\n", cmd, *dbPath)
	if err := sqlite.RunMigrations(db.DB, migrationsFS, "migrations", cmd); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration command completed successfully.")
}
