package balances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres/postgreserr"
	"github.com/samar/sup_bot/metacore/storage"
)

// BalanceStorage реализует интерфейс BalanceStorage.
type BalanceStorage struct {
	db storage.DBInterface
}

// NewBalanceStorage создает новый экземпляр BalanceStorage.
func NewBalanceStorage(db storage.DBInterface) *BalanceStorage {
	return &BalanceStorage{db: db}
}

// UpdateBalance обновляет баланс пользователя атомарно.
// Возвращает applied=true если баланс был обновлен (значения изменились).
func (s *BalanceStorage) UpdateBalance(ctx context.Context, balance *domain.UserBalance) (applied bool, err error) {
	query := `
		INSERT INTO user_balances (user_id, asset, free, locked, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, asset) 
		DO UPDATE SET 
			free = EXCLUDED.free,
			locked = EXCLUDED.locked,
			updated_at = EXCLUDED.updated_at
		WHERE user_balances.free != EXCLUDED.free 
		   OR user_balances.locked != EXCLUDED.locked
		RETURNING id`

	balance.UpdatedAt = time.Now()

	var id uint64
	err = s.db.QueryRowContext(ctx, query,
		balance.UserID,
		balance.Asset,
		balance.Free,
		balance.Locked,
		balance.UpdatedAt,
	).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Значения не изменились
			return false, nil
		}
		return false, fmt.Errorf("failed to update balance: %w", err)
	}

	balance.ID = id
	return true, nil
}

// GetBalance получает баланс пользователя по активу.
func (s *BalanceStorage) GetBalance(ctx context.Context, userID uint64, asset string) (*domain.UserBalance, error) {
	query := `
		SELECT id, user_id, asset, free, locked, updated_at
		FROM user_balances
		WHERE user_id = $1 AND asset = $2`

	balance := &domain.UserBalance{}
	err := s.db.QueryRowContext(ctx, query, userID, asset).Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Asset,
		&balance.Free,
		&balance.Locked,
		&balance.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, postgreserr.ErrBalanceNotFound
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// GetUserBalances получает все балансы пользователя.
func (s *BalanceStorage) GetUserBalances(ctx context.Context, userID uint64) ([]*domain.UserBalance, error) {
	query := `
		SELECT id, user_id, asset, free, locked, updated_at
		FROM user_balances
		WHERE user_id = $1
		ORDER BY asset`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user balances: %w", err)
	}
	defer rows.Close()

	var balances []*domain.UserBalance
	for rows.Next() {
		balance := &domain.UserBalance{}
		err := rows.Scan(
			&balance.ID,
			&balance.UserID,
			&balance.Asset,
			&balance.Free,
			&balance.Locked,
			&balance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan balance: %w", err)
		}
		balances = append(balances, balance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating balance rows: %w", err)
	}

	return balances, nil
}

// UpdateUserBalances обновляет все балансы пользователя атомарно в транзакции.
func (s *BalanceStorage) UpdateUserBalances(ctx context.Context, userID uint64, balances []*domain.UserBalance) error {
	if len(balances) == 0 {
		return nil
	}

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Подготавливаем запрос для обновления
	updateQuery := `
		INSERT INTO user_balances (user_id, asset, free, locked, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, asset) 
		DO UPDATE SET 
			free = EXCLUDED.free,
			locked = EXCLUDED.locked,
			updated_at = EXCLUDED.updated_at
		RETURNING id`

	stmt, err := tx.PrepareContext(ctx, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	updateTime := time.Now()

	// Обновляем каждый баланс
	for _, balance := range balances {
		balance.UserID = userID
		balance.UpdatedAt = updateTime

		var id uint64
		err = stmt.QueryRowContext(ctx,
			balance.UserID,
			balance.Asset,
			balance.Free,
			balance.Locked,
			balance.UpdatedAt,
		).Scan(&id)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to update balance for asset %s: %w", balance.Asset, err)
		}

		if err == nil {
			balance.ID = id
		}
	}

	// Удаляем балансы с нулевыми значениями
	deleteQuery := `
		DELETE FROM user_balances 
		WHERE user_id = $1 
		  AND free = 0 
		  AND locked = 0
		  AND updated_at < $2`

	_, err = tx.ExecContext(ctx, deleteQuery, userID, updateTime)
	if err != nil {
		return fmt.Errorf("failed to delete zero balances: %w", err)
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Ensure BalanceStorage implements BalanceStorage interface
var _ storage.BalanceStorage = (*BalanceStorage)(nil)
