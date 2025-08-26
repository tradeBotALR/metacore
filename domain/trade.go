package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trade представляет выполненную сделку в системе MEXC.
// Поля соответствуют таблице trades в БД и API MEXC.
type Trade struct {
	ID              uint64          `db:"id"`
	UserID          uint64          `db:"user_id"`
	MexcTradeID     string          `db:"mexc_trade_id"`
	OrderID         string          `db:"order_id"`
	Symbol          string          `db:"symbol"`
	Price           decimal.Decimal `db:"price"`
	Quantity        decimal.Decimal `db:"quantity"`
	QuoteQuantity   decimal.Decimal `db:"quote_quantity"`
	Commission      decimal.Decimal `db:"commission"`
	CommissionAsset string          `db:"commission_asset"`
	TradeTime       time.Time       `db:"trade_time"`
	IsBuyer         bool            `db:"is_buyer"`
	IsMaker         bool            `db:"is_maker"`
	CreatedAt       time.Time       `db:"created_at"`
}
