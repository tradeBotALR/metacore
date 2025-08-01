// cmd/deploy/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, *dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer func() {
		if closeErr := conn.Close(ctx); closeErr != nil {
			log.Printf("Warning: error closing database connection: %v", closeErr)
		}
	}()

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println("Connected to the database successfully.")

	//nolint:gosec
	if err := DeploySchema(ctx, conn, *schemaPath); err != nil {
		log.Fatalf("Schema deployment failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Deployment completed.")
}

func DeploySchema(ctx context.Context, conn *pgx.Conn, schemaPath string) error {
	//nolint:gosec
	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	sql := string(sqlBytes)

	_, err = conn.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Schema deployed successfully!")

	return nil
}
