package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"metacore/configs"
)

// DB represents a connection pool to the PostgreSQL database
type DB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(cfg configs.Config) (*DB, error) {
	// Формируем строку подключения из DBConfig
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode,
	)

	// Парсим конфигурацию пула
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	// Применяем настройки пула из PoolConfig
	config.MaxConns = cfg.Pool.MaxConns
	config.MinConns = cfg.Pool.MinConns
	config.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	config.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	config.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	ctx, cancel := context.WithTimeout(context.Background(), configs.DefaultConnectTimeout)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &DB{
		pool: db,
	}, nil
}

// Pool возвращает внутренний пул соединений.
// Используется для передачи в слои хранения (storage).
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
		log.Println("PostgreSQL connection pool closed")
	}
}

// Ping checks the database connection
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
