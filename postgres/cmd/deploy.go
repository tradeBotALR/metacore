// cmd/deploy/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const dbConnectTimeout = 10 * time.Second

var (
	dbURL      = flag.String("db-url", "", "PostgreSQL connection URL (e.g., postgres://user:pass@host:port/dbname)")
	schemaPath = flag.String("schema", "./postgres/schema.sql", "Path to the schema.sql file")
)

func main() {
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("Please provide a database URL using -db-url flag")
	}

	// 1. Подключаемся к БД
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()

	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err) //nolint:gocritic
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Warning: error closing database connection: %v", closeErr)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println("Connected to the database successfully.")

	//nolint:gosec
	if err := DeploySchema(ctx, db, *schemaPath); err != nil {
		log.Printf("Unable to deploy schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Deployment completed.")
}

func DeploySchema(ctx context.Context, db *sql.DB, schemaPath string) error {
	//nolint:gosec
	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	sql := string(sqlBytes)

	_, err = db.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Schema deployed successfully!")

	return nil
}
