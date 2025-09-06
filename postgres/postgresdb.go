package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/samar/sup_bot/metacore/configs"
	"github.com/samar/sup_bot/metacore/postgres/internal/balances"
	"github.com/samar/sup_bot/metacore/postgres/internal/orders"
	"github.com/samar/sup_bot/metacore/postgres/internal/trades"
	"github.com/samar/sup_bot/metacore/postgres/internal/users"
	"github.com/samar/sup_bot/metacore/storage"

	_ "github.com/lib/pq"
)

// DB represents a connection to the PostgreSQL database
type DB struct {
	db *sql.DB
	storage.FullStorage
	UserStorage storage.UserStorage
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg configs.Config) (*DB, error) {
	// Формируем строку подключения из DBConfig
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode,
	)

	// Создаем соединение с БД
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(int(cfg.Pool.MaxConns))
	db.SetMaxIdleConns(int(cfg.Pool.MinConns))
	db.SetConnMaxLifetime(cfg.Pool.MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.Pool.MaxConnIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), configs.DefaultConnectTimeout)
	defer cancel()

	// Проверяем соединение
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Создаем storage слои
	dbAdapter := storage.NewDBAdapter(db)
	orderStorage := orders.NewOrderStorage(dbAdapter)
	userStorage := users.NewUserStorage(dbAdapter)
	tradeStorage := trades.NewTradeStorage(dbAdapter)
	balanceStorage := balances.NewBalanceStorage(dbAdapter)
	orderUpdateStorage := orders.NewOrderUpdateStorage(dbAdapter)

	// Создаем FullStorage, объединяющий все storage
	fullStorage := &fullStorage{
		UserStorage:        userStorage,
		OrderStorage:       orderStorage,
		TradeStorage:       tradeStorage,
		BalanceStorage:     balanceStorage,
		OrderUpdateStorage: orderUpdateStorage,
	}

	return &DB{
		db:          db,
		FullStorage: fullStorage,
		UserStorage: userStorage,
	}, nil
}

// Close закрывает соединение с БД
func (db *DB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

// GetUserStorage возвращает хранилище пользователей
func (db *DB) GetUserStorage() storage.UserStorage {
	return db.UserStorage
}

// Ping проверяет соединение
func (db *DB) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// fullStorage объединяет все storage интерфейсы
type fullStorage struct {
	storage.UserStorage
	storage.OrderStorage
	storage.TradeStorage
	storage.BalanceStorage
	storage.OrderUpdateStorage
}

// Ensure fullStorage implements FullStorage interface
var _ storage.FullStorage = (*fullStorage)(nil)
