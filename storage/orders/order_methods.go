package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"metacore/domain"
	"metacore/postgres/postgreserr"
	"metacore/storage"
)

// Storage реализует интерфейс DBStorage.
type Storage struct {
	db *pgxpool.Pool
}

// NewStorage создает новый экземпляр Storage.
func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

// CreateOrder сохраняет новый ордер в хранилище.
func (s *Storage) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
        INSERT INTO orders (
            internal_id, user_id, mexc_order_id, symbol, side, type, status,
            price, quantity, quote_order_qty, executed_quantity,
            cummulative_quote_qty, client_order_id, transact_time
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )`

	_, err := s.db.Exec(ctx, query,
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
func (s *Storage) DeleteOrderByID(ctx context.Context, mexcOrderID string) error {
	query := `DELETE FROM orders WHERE mexc_order_id = $1`

	result, err := s.db.Exec(ctx, query, mexcOrderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	// Проверка, была ли удалена хотя бы одна строка
	if result.RowsAffected() == 0 {
		// Можно вернуть кастомную ошибку, если это важно
		return fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
	}

	return nil
}

// UpdateOrderStatus обновляет статус ордера.
func (s *Storage) UpdateOrderStatus(ctx context.Context, mexcOrderID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE mexc_order_id = $2`

	result, err := s.db.Exec(ctx, query, status, mexcOrderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
	}

	return nil
}

// GetOrderByID получает ордер по его mexc_order_id.
func (s *Storage) GetOrderByID(ctx context.Context, mexcOrderID string) (*domain.Order, error) {
	query := `
        SELECT id, internal_id, user_id, mexc_order_id, symbol, side, type, status,
               price, quantity, quote_order_qty, executed_quantity,
               cummulative_quote_qty, client_order_id, transact_time, -- Возвращаем в миллисекундах
               created_at,updated_at
        FROM orders WHERE mexc_order_id = $1`

	var order domain.Order

	err := s.db.QueryRow(ctx, query, mexcOrderID).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			// Можно вернуть кастомную ошибку
			return nil, fmt.Errorf("order with id %s not found: %w", mexcOrderID, postgreserr.ErrOrderNotFound)
		}

		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// Close закрывает соединение с БД.
func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ensure Storage implements DBStorage interface
var _ storage.DBStorage = (*Storage)(nil)
