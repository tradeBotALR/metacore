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
			telegram_id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
			kyc_status, can_trade, can_withdraw, can_deposit,
			account_type, permissions, last_account_sync, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query,
		user.TelegramID,
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
		SELECT id, telegram_id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users WHERE id = $1`

	var user domain.User

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.TelegramID,
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
		SELECT id, telegram_id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users WHERE mexc_uid = $1`

	var user domain.User

	err := s.db.QueryRowContext(ctx, query, mexcUID).Scan(
		&user.ID,
		&user.TelegramID,
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

// GetUserByTelegramID получает пользователя по Telegram ID.
func (s *UserStorage) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `
		SELECT id, telegram_id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users WHERE telegram_id = $1`

	var user domain.User

	err := s.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
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
			return nil, fmt.Errorf("user with telegram_id %d not found: %w", telegramID, postgreserr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser обновляет пользователя.
func (s *UserStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET 
			telegram_id = $1, mexc_uid = $2, username = $3, email = $4, mexc_api_key = $5,
			mexc_secret_key = $6, kyc_status = $7, can_trade = $8,
			can_withdraw = $9, can_deposit = $10, account_type = $11,
			permissions = $12, last_account_sync = $13, is_active = $14,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $15`

	result, err := s.db.ExecContext(ctx, query,
		user.TelegramID,
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

// GetAllUsers получает всех пользователей.
func (s *UserStorage) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, telegram_id, mexc_uid, username, email, mexc_api_key, mexc_secret_key,
		       kyc_status, can_trade, can_withdraw, can_deposit, account_type,
		       permissions, last_account_sync, is_active, created_at, updated_at
		FROM users ORDER BY id`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.TelegramID,
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
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}

// Close закрывает соединение с БД.
func (s *UserStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ensure UserStorage implements UserStorage interface
var _ storage.UserStorage = (*UserStorage)(nil)
