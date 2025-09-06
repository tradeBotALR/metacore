package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

// GetUserOrders получает ордера пользователя с фильтрами
func (s *OrderStorage) GetUserOrders(ctx context.Context, userID uint64, filters ...storage.OrderFilter) ([]*domain.Order, error) {
	var filter storage.OrderFilter
	if len(filters) > 0 {
		filter = filters[0]
	}

	b := &strings.Builder{}
	b.WriteString(`SELECT id, internal_id, user_id, mexc_order_id, symbol, side, type, status,
price, quantity, quote_order_qty, executed_quantity, cummulative_quote_qty, client_order_id, transact_time, created_at, updated_at
FROM orders WHERE user_id = $1`)
	args := []interface{}{userID}
	idx := 1

	if filter.Symbol != "" {
		idx++
		b.WriteString(fmt.Sprintf(" AND symbol = $%d", idx))
		args = append(args, filter.Symbol)
	}
	if filter.Status != "" {
		idx++
		b.WriteString(fmt.Sprintf(" AND status = $%d", idx))
		args = append(args, filter.Status)
	}
	if filter.StartTime != nil {
		idx++
		b.WriteString(fmt.Sprintf(" AND transact_time >= $%d", idx))
		args = append(args, *filter.StartTime)
	}
	if filter.EndTime != nil {
		idx++
		b.WriteString(fmt.Sprintf(" AND transact_time <= $%d", idx))
		args = append(args, *filter.EndTime)
	}
	b.WriteString(" ORDER BY created_at DESC, id DESC")
	if filter.Limit > 0 {
		idx++
		b.WriteString(fmt.Sprintf(" LIMIT $%d", idx))
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		idx++
		b.WriteString(fmt.Sprintf(" OFFSET $%d", idx))
		args = append(args, filter.Offset)
	}

	rows, err := s.db.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.InternalID, &o.UserID, &o.MexcOrderID, &o.Symbol, &o.Side, &o.Type, &o.Status,
			&o.Price, &o.Quantity, &o.QuoteOrderQty, &o.ExecutedQuantity, &o.CummulativeQuoteQty, &o.ClientOrderID,
			&o.TransactTime, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("orders rows error: %w", err)
	}
	return orders, nil
}

// GetOpenOrders получает активные ордера пользователя
func (s *OrderStorage) GetOpenOrders(ctx context.Context, userID uint64, symbol string) ([]*domain.Order, error) {
	b := &strings.Builder{}
	b.WriteString(`SELECT id, internal_id, user_id, mexc_order_id, symbol, side, type, status,
price, quantity, quote_order_qty, executed_quantity, cummulative_quote_qty, client_order_id, transact_time, created_at, updated_at
FROM orders WHERE user_id = $1 AND status IN ('NEW','PARTIALLY_FILLED')`)
	args := []interface{}{userID}
	if symbol != "" {
		b.WriteString(" AND symbol = $2")
		args = append(args, symbol)
	}
	b.WriteString(" ORDER BY created_at DESC, id DESC")

	rows, err := s.db.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query open orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.InternalID, &o.UserID, &o.MexcOrderID, &o.Symbol, &o.Side, &o.Type, &o.Status,
			&o.Price, &o.Quantity, &o.QuoteOrderQty, &o.ExecutedQuantity, &o.CummulativeQuoteQty, &o.ClientOrderID,
			&o.TransactTime, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan open order: %w", err)
		}
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("open orders rows error: %w", err)
	}
	return orders, nil
}

// --- Order updates ---

type OrderUpdateStorageImpl struct {
	db storage.DBInterface
}

func NewOrderUpdateStorage(db storage.DBInterface) *OrderUpdateStorageImpl {
	return &OrderUpdateStorageImpl{db: db}
}

func (s *OrderUpdateStorageImpl) AppendOrderUpdate(ctx context.Context, update *domain.OrderUpdate) error {
	query := `INSERT INTO order_updates (user_id, order_id, status, executed_quantity, cummulative_quote_qty, update_time, raw_data)
VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := s.db.ExecContext(ctx, query,
		update.UserID, update.OrderID, update.Status,
		update.ExecutedQuantity, update.CummulativeQuoteQty,
		update.UpdateTime, update.RawData,
	)
	if err != nil {
		return fmt.Errorf("failed to append order update: %w", err)
	}
	return nil
}

func (s *OrderUpdateStorageImpl) GetOrderUpdates(ctx context.Context, userID uint64, orderID string) ([]*domain.OrderUpdate, error) {
	query := `SELECT id, user_id, order_id, status, executed_quantity, cummulative_quote_qty, update_time, raw_data
FROM order_updates WHERE user_id = $1 AND order_id = $2 ORDER BY update_time DESC, id DESC`
	rows, err := s.db.QueryContext(ctx, query, userID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order updates: %w", err)
	}
	defer rows.Close()

	var updates []*domain.OrderUpdate
	for rows.Next() {
		u := &domain.OrderUpdate{}
		if err := rows.Scan(&u.ID, &u.UserID, &u.OrderID, &u.Status, &u.ExecutedQuantity, &u.CummulativeQuoteQty, &u.UpdateTime, &u.RawData); err != nil {
			return nil, fmt.Errorf("failed to scan order update: %w", err)
		}
		updates = append(updates, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("order updates rows error: %w", err)
	}
	return updates, nil
}

// Close закрывает соединение с БД.
func (s *OrderStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ensure OrderStorage implements OrderStorage interface
var _ storage.OrderStorage = (*OrderStorage)(nil)
