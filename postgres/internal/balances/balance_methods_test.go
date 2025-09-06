package balances

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/storage"
)

func TestBalanceStorage_UpdateBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update with changes", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		balance := &domain.UserBalance{
			UserID: 1,
			Asset:  "BTC",
			Free:   decimal.NewFromFloat(1.5),
			Locked: decimal.NewFromFloat(0.5),
		}

		mock.ExpectQuery("INSERT INTO user_balances").
			WithArgs(balance.UserID, balance.Asset, balance.Free, balance.Locked, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

		applied, err := s.UpdateBalance(ctx, balance)
		assert.NoError(t, err)
		assert.True(t, applied)
		assert.Equal(t, uint64(10), balance.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no changes made", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		balance := &domain.UserBalance{
			UserID: 1,
			Asset:  "BTC",
			Free:   decimal.NewFromFloat(1.5),
			Locked: decimal.NewFromFloat(0.5),
		}

		mock.ExpectQuery("INSERT INTO user_balances").
			WithArgs(balance.UserID, balance.Asset, balance.Free, balance.Locked, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		applied, err := s.UpdateBalance(ctx, balance)
		assert.NoError(t, err)
		assert.False(t, applied)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		balance := &domain.UserBalance{
			UserID: 1,
			Asset:  "BTC",
			Free:   decimal.NewFromFloat(1.5),
			Locked: decimal.NewFromFloat(0.5),
		}

		mock.ExpectQuery("INSERT INTO user_balances").
			WithArgs(balance.UserID, balance.Asset, balance.Free, balance.Locked, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		applied, err := s.UpdateBalance(ctx, balance)
		assert.Error(t, err)
		assert.False(t, applied)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBalanceStorage_GetUserBalances(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get multiple balances", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)
		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "user_id", "asset", "free", "locked", "updated_at"}).
			AddRow(1, userID, "BTC", decimal.NewFromFloat(1.5), decimal.NewFromFloat(0.5), now).
			AddRow(2, userID, "USDT", decimal.NewFromFloat(1000), decimal.NewFromFloat(0), now)

		mock.ExpectQuery("SELECT id, user_id, asset, free, locked, updated_at FROM user_balances").
			WithArgs(userID).
			WillReturnRows(rows)

		balances, err := s.GetUserBalances(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, balances, 2)
		assert.Equal(t, "BTC", balances[0].Asset)
		assert.Equal(t, "USDT", balances[1].Asset)
		assert.Equal(t, "1.5", balances[0].Free.String())
		assert.Equal(t, "1000", balances[1].Free.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty balances", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)

		rows := sqlmock.NewRows([]string{"id", "user_id", "asset", "free", "locked", "updated_at"})

		mock.ExpectQuery("SELECT id, user_id, asset, free, locked, updated_at FROM user_balances").
			WithArgs(userID).
			WillReturnRows(rows)

		balances, err := s.GetUserBalances(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, balances, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)

		mock.ExpectQuery("SELECT id, user_id, asset, free, locked, updated_at FROM user_balances").
			WithArgs(userID).
			WillReturnError(errors.New("query error"))

		balances, err := s.GetUserBalances(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, balances)
		assert.Contains(t, err.Error(), "query error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBalanceStorage_UpdateUserBalances(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update multiple balances", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)
		balances := []*domain.UserBalance{
			{
				Asset:  "BTC",
				Free:   decimal.NewFromFloat(1.5),
				Locked: decimal.NewFromFloat(0.5),
			},
			{
				Asset:  "USDT",
				Free:   decimal.NewFromFloat(1000),
				Locked: decimal.NewFromFloat(0),
			},
		}

		mock.ExpectBegin()

		// Prepare statement
		mock.ExpectPrepare("INSERT INTO user_balances")

		// First balance update
		mock.ExpectQuery("INSERT INTO user_balances").
			WithArgs(userID, "BTC", balances[0].Free, balances[0].Locked, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Second balance update
		mock.ExpectQuery("INSERT INTO user_balances").
			WithArgs(userID, "USDT", balances[1].Free, balances[1].Locked, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		// Delete zero balances
		mock.ExpectExec("DELETE FROM user_balances").
			WithArgs(userID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		err = s.UpdateUserBalances(ctx, userID, balances)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), balances[0].ID)
		assert.Equal(t, uint64(2), balances[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty balances list", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)

		err = s.UpdateUserBalances(ctx, userID, []*domain.UserBalance{})
		assert.NoError(t, err)
	})

	t.Run("begin transaction error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)
		balances := []*domain.UserBalance{
			{Asset: "BTC", Free: decimal.NewFromFloat(1.5), Locked: decimal.NewFromFloat(0.5)},
		}

		mock.ExpectBegin().WillReturnError(errors.New("tx error"))

		err = s.UpdateUserBalances(ctx, userID, balances)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tx error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("prepare statement error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewBalanceStorage(dbAdapter)

		userID := uint64(1)
		balances := []*domain.UserBalance{
			{Asset: "BTC", Free: decimal.NewFromFloat(1.5), Locked: decimal.NewFromFloat(0.5)},
		}

		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO user_balances").WillReturnError(errors.New("prepare error"))
		mock.ExpectRollback()

		err = s.UpdateUserBalances(ctx, userID, balances)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "prepare error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
