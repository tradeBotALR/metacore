package storage

import (
	"context"
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

// DBStorage объединяет все интерфейсы хранилища.
// В дальнейшем сюда можно добавить UserStorage, TradeStorage и т.д.
type DBStorage interface {
	OrderStorage
	// Close закрывает соединение с БД.
	Close()
}
