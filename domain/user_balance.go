package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserBalance представляет баланс пользователя по конкретному активу.
// Поля соответствуют таблице user_balances в БД.
type UserBalance struct {
	ID        uint64          `db:"id"`
	UserID    uint64          `db:"user_id"`
	Asset     string          `db:"asset"`
	Free      decimal.Decimal `db:"free"`   // Доступный баланс
	Locked    decimal.Decimal `db:"locked"` // Заблокированный баланс
	UpdatedAt time.Time       `db:"updated_at"`
}
