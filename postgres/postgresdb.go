package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnectTimeout    = 10 * time.Second // Время ожидания подключения
	defaultMaxConnIdleTime   = 30 * time.Minute // Максимальное время простоя соединения
	defaultMaxConnLifetime   = 1 * time.Hour    // Максимальное время жизни соединения
	defaultHealthCheckPeriod = 1 * time.Minute  // Период проверки состояния соединений
	defaultMaxConns          = int32(20)        // Максимальное количество соединений в пуле
	defaultMinConns          = int32(5)         // Минимальное количество соединений в пуле
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DB represents a connection pool to the PostgreSQL database
type DB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(cfg Config) (*DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Configure the connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	// Set pool configuration using constants
	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.MaxConnLifetime = defaultMaxConnLifetime
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.HealthCheckPeriod = defaultHealthCheckPeriod

	// Attempt to connect to the database
	// mnd: Magic number - заменено на константу
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	// wrapcheck: Оборачиваем ошибку из внешнего пакета
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &DB{
		pool: db,
	}, nil
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
