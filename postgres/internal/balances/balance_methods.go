package balances

import (
	"context"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/storage"
)

// BalanceStorage реализует интерфейс BalanceStorage.
type BalanceStorage struct {
	// TODO: Реализовать методы для работы с балансами
}

// NewBalanceStorage создает новый экземпляр BalanceStorage.
func NewBalanceStorage() *BalanceStorage {
	return &BalanceStorage{}
}

// UpdateBalance обновляет баланс пользователя.
func (s *BalanceStorage) UpdateBalance(ctx context.Context, balance *domain.UserBalance) error {
	// TODO: Реализовать обновление баланса
	return nil
}

// GetBalance получает баланс пользователя по активу.
func (s *BalanceStorage) GetBalance(ctx context.Context, userID uint64, asset string) (*domain.UserBalance, error) {
	// TODO: Реализовать получение баланса
	return nil, nil
}

// Ensure BalanceStorage implements BalanceStorage interface
var _ storage.BalanceStorage = (*BalanceStorage)(nil)
