package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres/postgreserr"
	"github.com/samar/sup_bot/metacore/storage"
)

// OrderStorage реализует интерфейс OrderStorage.
type OrderStorage struct {
	db storage.DBInterface
}

// NewOrderStorage создает новый экземпляр OrderStorage.
func NewOrderStorage(db storage.DBInterface) *OrderStorage {
	return &OrderStorage{db: db}
}

// CreateOrder сохраняет новый ордер в хранилище.
func (s *OrderStorage) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
        INSERT INTO orders (
            internal_id, user_id, mexc_order_id, symbol, side, type, status,
            price, quantity, quote_order_qty, executed_quantity,
            cummulative_quote_qty, client_order_id, transact_time
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )`

	_, err := s.db.ExecContext(ctx, query,
		order.InternalID,
		order.UserID,
		order.MexcOrderID,
		order.Symbol,
		order.Side,
		order.Type,
		order.Status,
		order.Price,
		order.Quantity,
		order.QuoteOrderQty,
		order.ExecutedQuantity,
		order.CummulativeQuoteQty,
		order.ClientOrderID,
		order.TransactTime,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// DeleteOrderByID удаляет ордер из хранилища по его mexc_order_id.
func (s *OrderStorage) DeleteOrderByID(ctx context.Context, mexcOrderID string) error {
	query := `DELETE FROM orders WHERE mexc_order_id = $1`

	result, err := s.db.ExecContext(ctx, query, mexcOrderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	// Проверка, была ли удалена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Можно вернуть кастомную ошибку, если это важно
		return fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
	}

	return nil
}

// UpdateOrderStatus обновляет статус ордера.
func (s *OrderStorage) UpdateOrderStatus(ctx context.Context, mexcOrderID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE mexc_order_id = $2`

	result, err := s.db.ExecContext(ctx, query, status, mexcOrderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
	}

	return nil
}

// GetOrderByID получает ордер по его mexc_order_id.
func (s *OrderStorage) GetOrderByID(ctx context.Context, mexcOrderID string) (*domain.Order, error) {
	query := `
        SELECT id, internal_id, user_id, mexc_order_id, symbol, side, type, status,
               price, quantity, quote_order_qty, executed_quantity,
               cummulative_quote_qty, client_order_id, transact_time, -- Возвращаем в миллисекундах
               created_at,updated_at
        FROM orders WHERE mexc_order_id = $1`

	var order domain.Order

	err := s.db.QueryRowContext(ctx, query, mexcOrderID).Scan(
		&order.ID,
		&order.InternalID,
		&order.UserID,
		&order.MexcOrderID,
		&order.Symbol,
		&order.Side,
		&order.Type,
		&order.Status,
		&order.Price,
		&order.Quantity,
		&order.QuoteOrderQty,
		&order.ExecutedQuantity,
		&order.CummulativeQuoteQty,
		&order.ClientOrderID,
		&order.TransactTime,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Можно вернуть кастомную ошибку
			return nil, fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
		}

		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// Close закрывает соединение с БД.
func (s *OrderStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ensure OrderStorage implements OrderStorage interface
var _ storage.OrderStorage = (*OrderStorage)(nil)
