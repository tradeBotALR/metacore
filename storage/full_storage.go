package storage

import (
	"context"
	"database/sql"

	"github.com/samar/sup_bot/metacore/domain"
)

// RowInterface представляет интерфейс для sql.Row
type RowInterface interface {
	Scan(dest ...interface{}) error
}

type OrderStorage interface {
	// CreateOrder сохраняет новый ордер в хранилище.
	CreateOrder(ctx context.Context, order *domain.Order) error

	// DeleteOrderByID удаляет ордер из хранилища по его mexc_order_id.
	DeleteOrderByID(ctx context.Context, mexcOrderID string) error

	// UpdateOrderStatus обновляет статус ордера (полезно будет сразу)
	UpdateOrderStatus(ctx context.Context, mexcOrderID, status string) error

	// GetOrderByID получает ордер по его mexc_order_id (полезно будет сразу)
	GetOrderByID(ctx context.Context, mexcOrderID string) (*domain.Order, error)
}

type UserStorage interface {
	// CreateUser создает нового пользователя
	CreateUser(ctx context.Context, user *domain.User) error

	// GetUserByID получает пользователя по ID
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)

	// GetUserByMexcUID получает пользователя по MEXC UID
	GetUserByMexcUID(ctx context.Context, mexcUID string) (*domain.User, error)

	// UpdateUser обновляет пользователя
	UpdateUser(ctx context.Context, user *domain.User) error

	// DeleteUser удаляет пользователя
	DeleteUser(ctx context.Context, id uint64) error
}

type TradeStorage interface {
	// CreateTrade создает новую сделку
	CreateTrade(ctx context.Context, trade *domain.Trade) error

	// GetTradeByID получает сделку по MEXC Trade ID
	GetTradeByID(ctx context.Context, mexcTradeID string) (*domain.Trade, error)
}

type BalanceStorage interface {
	// UpdateBalance обновляет баланс пользователя
	UpdateBalance(ctx context.Context, balance *domain.UserBalance) error

	// GetBalance получает баланс пользователя по активу
	GetBalance(ctx context.Context, userID uint64, asset string) (*domain.UserBalance, error)
}

// FullStorage объединяет все интерфейсы хранилища.
type FullStorage interface {
	UserStorage
	OrderStorage
	TradeStorage
	BalanceStorage
}

// DBInterface определяет интерфейс для работы с базой данных,
// который будет реализован *sql.DB.
// Это позволяет легко мокать его в unit-тестах.
type DBInterface interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) RowInterface
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	PingContext(ctx context.Context) error
	Close() error
}
