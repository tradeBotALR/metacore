package domain

import (
	"github.com/shopspring/decimal"
	"time"
) // Для точной работы с финансами

// Order представляет собой ордер в системе.
// Поля должны соответствовать таблице orders в БД.
type Order struct {
	ID                  uint64          `db:"id"`
	InternalID          int64           `db:"internal_id"` // BIGINT UNIQUE
	UserID              uint64          `db:"user_id"`
	MexcOrderID         string          `db:"mexc_order_id"`
	Symbol              string          `db:"symbol"`
	Side                string          `db:"side"` // BUY, SELL
	Type                string          `db:"type"` // LIMIT, MARKET, etc.
	Status              string          `db:"status"`
	Price               decimal.Decimal `db:"price"` // DECIMAL(30, 15)
	Quantity            decimal.Decimal `db:"quantity"`
	QuoteOrderQty       decimal.Decimal `db:"quote_order_qty"` // Может быть NULL
	ExecutedQuantity    decimal.Decimal `db:"executed_quantity"`
	CummulativeQuoteQty decimal.Decimal `db:"cummulative_quote_qty"`
	ClientOrderID       string          `db:"client_order_id"` // Может быть NULL
	TransactTime        time.Time       `db:"transact_time"`   // Unix timestamp в миллисекундах из API
	CreatedAt           time.Time       `db:"created_at"`      // Можно добавить, если нужно в коде
	UpdatedAt           time.Time       `db:"updated_at"`      // Можно добавить, если нужно в коде
}
