// Package metacore предоставляет основную функциональность для работы с базой данных,
// пользователями, ордерами, сделками и балансами.
package metacore

import (
	"github.com/samar/sup_bot/metacore/configs"
	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres"
	"github.com/samar/sup_bot/metacore/storage"
)

// NewPostgresDB создает новое подключение к PostgreSQL базе данных
func NewPostgresDB(cfg configs.Config) (*postgres.DB, error) {
	return postgres.NewPostgresDB(cfg)
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() configs.Config {
	return configs.DefaultConfig()
}

// User представляет пользователя в системе
type User = domain.User

// Order представляет ордер в системе
type Order = domain.Order

// Trade представляет сделку в системе
type Trade = domain.Trade

// UserBalance представляет баланс пользователя
type UserBalance = domain.UserBalance

// FullStorage интерфейс для работы со всеми типами хранилищ
type FullStorage = storage.FullStorage

// DBInterface интерфейс для работы с базой данных
type DBInterface = storage.DBInterface

// RowInterface интерфейс для работы с результатами запросов
type RowInterface = storage.RowInterface

// OrderStorage интерфейс для работы с ордерами
type OrderStorage = storage.OrderStorage

// UserStorage интерфейс для работы с пользователями
type UserStorage = storage.UserStorage

// TradeStorage интерфейс для работы со сделками
type TradeStorage = storage.TradeStorage

// BalanceStorage интерфейс для работы с балансами
type BalanceStorage = storage.BalanceStorage
