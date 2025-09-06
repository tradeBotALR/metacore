package trades

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

// Re-export TradeFilter for convenience
type TradeFilter = storage.TradeFilter

// TradeStorage реализует интерфейс TradeStorage.
type TradeStorage struct {
	db storage.DBInterface
}

// NewTradeStorage создает новый экземпляр TradeStorage.
func NewTradeStorage(db storage.DBInterface) *TradeStorage {
	return &TradeStorage{db: db}
}

// CreateTrade создает новую сделку.
func (s *TradeStorage) CreateTrade(ctx context.Context, trade *domain.Trade) error {
	query := `
		INSERT INTO trades (
			user_id, mexc_trade_id, order_id, symbol, price, quantity,
			quote_quantity, commission, commission_asset, trade_time,
			is_buyer, is_maker
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query,
		trade.UserID,
		trade.MexcTradeID,
		trade.OrderID,
		trade.Symbol,
		trade.Price,
		trade.Quantity,
		trade.QuoteQuantity,
		trade.Commission,
		trade.CommissionAsset,
		trade.TradeTime,
		trade.IsBuyer,
		trade.IsMaker,
	).Scan(&trade.ID, &trade.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	return nil
}

// GetTradeByID получает сделку по MEXC Trade ID.
func (s *TradeStorage) GetTradeByID(ctx context.Context, mexcTradeID string) (*domain.Trade, error) {
	query := `
		SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity,
			   quote_quantity, commission, commission_asset, trade_time,
			   is_buyer, is_maker, created_at
		FROM trades
		WHERE mexc_trade_id = $1`

	trade := &domain.Trade{}
	err := s.db.QueryRowContext(ctx, query, mexcTradeID).Scan(
		&trade.ID,
		&trade.UserID,
		&trade.MexcTradeID,
		&trade.OrderID,
		&trade.Symbol,
		&trade.Price,
		&trade.Quantity,
		&trade.QuoteQuantity,
		&trade.Commission,
		&trade.CommissionAsset,
		&trade.TradeTime,
		&trade.IsBuyer,
		&trade.IsMaker,
		&trade.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, postgreserr.ErrTradeNotFound
		}
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return trade, nil
}

// GetUserTrades получает все сделки пользователя с фильтрацией.
func (s *TradeStorage) GetUserTrades(ctx context.Context, userID uint64, filters ...storage.TradeFilter) ([]*domain.Trade, error) {
	var filter storage.TradeFilter
	if len(filters) > 0 {
		filter = filters[0]
	}

	// Базовый запрос
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity,
			   quote_quantity, commission, commission_asset, trade_time,
			   is_buyer, is_maker, created_at
		FROM trades
		WHERE user_id = $1`)

	args := []interface{}{userID}
	argCount := 1

	// Добавляем фильтры
	if filter.Symbol != "" {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" AND symbol = $%d", argCount))
		args = append(args, filter.Symbol)
	}

	if filter.StartTime != nil {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" AND trade_time >= $%d", argCount))
		args = append(args, *filter.StartTime)
	}

	if filter.EndTime != nil {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" AND trade_time <= $%d", argCount))
		args = append(args, *filter.EndTime)
	}

	// Сортировка - новые сделки первыми
	queryBuilder.WriteString(" ORDER BY trade_time DESC, id DESC")

	// Лимит и оффсет
	if filter.Limit > 0 {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argCount))
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argCount))
		args = append(args, filter.Offset)
	}

	query := queryBuilder.String()
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user trades: %w", err)
	}
	defer rows.Close()

	var trades []*domain.Trade
	for rows.Next() {
		trade := &domain.Trade{}
		err := rows.Scan(
			&trade.ID,
			&trade.UserID,
			&trade.MexcTradeID,
			&trade.OrderID,
			&trade.Symbol,
			&trade.Price,
			&trade.Quantity,
			&trade.QuoteQuantity,
			&trade.Commission,
			&trade.CommissionAsset,
			&trade.TradeTime,
			&trade.IsBuyer,
			&trade.IsMaker,
			&trade.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trade rows: %w", err)
	}

	return trades, nil
}

// Ensure TradeStorage implements TradeStorage interface
var _ storage.TradeStorage = (*TradeStorage)(nil)
