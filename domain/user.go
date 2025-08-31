package domain

import (
	"time"
)

// User представляет пользователя в системе MEXC.
// Поля соответствуют таблице users в БД и API MEXC.
type User struct {
	ID              uint64     `db:"id"`
	TelegramID      int64      `db:"telegram_id"`
	MexcUID         string     `db:"mexc_uid"`
	Username        string     `db:"username"`
	Email           string     `db:"email"`
	MexcAPIKey      string     `db:"mexc_api_key"`
	MexcSecretKey   string     `db:"mexc_secret_key"`
	KYCStatus       int16      `db:"kyc_status"`
	CanTrade        bool       `db:"can_trade"`
	CanWithdraw     bool       `db:"can_withdraw"`
	CanDeposit      bool       `db:"can_deposit"`
	AccountType     string     `db:"account_type"`
	Permissions     string     `db:"permissions"` // JSON array
	LastAccountSync *time.Time `db:"last_account_sync"`
	IsActive        bool       `db:"is_active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}
