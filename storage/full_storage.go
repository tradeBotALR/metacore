package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"metacore/domain"
)

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

// FullStorage объединяет все интерфейсы хранилища.
// В дальнейшем сюда можно добавить UserStorage, TradeStorage и т.д.
type FullStorage interface {
	OrderStorage
}

// PgxPoolIface определяет интерфейс для работы с пулом соединений,
// который будет реализован *pgxpool.Pool.
// Это позволяет легко мокать его в unit-тестах.
type PgxPoolIface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	// Добавь другие методы, если они будут использоваться
	// Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	// Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}
