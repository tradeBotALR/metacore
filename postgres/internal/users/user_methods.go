package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres/postgreserr"
	"github.com/samar/sup_bot/metacore/storage"
)

// UserStorage реализует интерфейс UserStorage.
type UserStorage struct {
	db storage.DBInterface
}

// NewUserStorage создает новый экземпляр UserStorage.
func NewUserStorage(db storage.DBInterface) *UserStorage {
	return &UserStorage{db: db}
}

// CreateUser создает нового пользователя.
func (s *UserStorage) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			mexc_uid, username, email, mexc_api_key, mexc_secret_key,
			kyc_status, can_trade, can_withdraw, can_deposit,
			account_type, permissions, last_account_sync, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query,
		user.MexcUID,
		user.Username,
		user.Email,
		user.MexcAPIKey,
		user.MexcSecretKey,
		user.KYCStatus,
		user.CanTrade,
		user.CanWithdraw,
		user.CanDeposit,
		user.AccountType,
		user.Permissions,
		user.LastAccountSync,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID получает пользователя по ID.
func (s *UserStorage) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	query := `
		SELECT id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users WHERE id = $1`

	var user domain.User

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.MexcUID,
		&user.Username,
		&user.Email,
		&user.MexcAPIKey,
		&user.MexcSecretKey,
		&user.KYCStatus,
		&user.CanTrade,
		&user.CanWithdraw,
		&user.CanDeposit,
		&user.AccountType,
		&user.Permissions,
		&user.LastAccountSync,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found: %w", id, postgreserr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByMexcUID получает пользователя по MEXC UID.
func (s *UserStorage) GetUserByMexcUID(ctx context.Context, mexcUID string) (*domain.User, error) {
	query := `
		SELECT id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users WHERE mexc_uid = $1`

	var user domain.User

	err := s.db.QueryRowContext(ctx, query, mexcUID).Scan(
		&user.ID,
		&user.MexcUID,
		&user.Username,
		&user.Email,
		&user.MexcAPIKey,
		&user.MexcSecretKey,
		&user.KYCStatus,
		&user.CanTrade,
		&user.CanWithdraw,
		&user.CanDeposit,
		&user.AccountType,
		&user.Permissions,
		&user.LastAccountSync,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with mexc_uid %s not found: %w", mexcUID, postgreserr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser обновляет пользователя.
func (s *UserStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET 
			mexc_uid = $1, username = $2, email = $3, mexc_api_key = $4,
			mexc_secret_key = $5, kyc_status = $6, can_trade = $7,
			can_withdraw = $8, can_deposit = $9, account_type = $10,
			permissions = $11, last_account_sync = $12, is_active = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $14`

	result, err := s.db.ExecContext(ctx, query,
		user.MexcUID,
		user.Username,
		user.Email,
		user.MexcAPIKey,
		user.MexcSecretKey,
		user.KYCStatus,
		user.CanTrade,
		user.CanWithdraw,
		user.CanDeposit,
		user.AccountType,
		user.Permissions,
		user.LastAccountSync,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found: %w", user.ID, postgreserr.ErrUserNotFound)
	}

	return nil
}

// DeleteUser удаляет пользователя.
func (s *UserStorage) DeleteUser(ctx context.Context, id uint64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found: %w", id, postgreserr.ErrUserNotFound)
	}

	return nil
}

// Close закрывает соединение с БД.
func (s *UserStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ensure UserStorage implements UserStorage interface
var _ storage.UserStorage = (*UserStorage)(nil)
