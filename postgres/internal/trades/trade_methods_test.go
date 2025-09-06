package trades

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
	"github.com/samar/sup_bot/metacore/postgres/postgreserr"
	"github.com/samar/sup_bot/metacore/storage"
)

func TestTradeStorage_CreateTrade(t *testing.T) {
	ctx := context.Background()

	t.Run("successful create", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		tradeTime := time.Now().Add(-time.Hour)
		createdAt := time.Now()
		trade := &domain.Trade{
			UserID:          1,
			MexcTradeID:     "123456",
			OrderID:         "order123",
			Symbol:          "BTCUSDT",
			Price:           decimal.NewFromFloat(50000),
			Quantity:        decimal.NewFromFloat(0.01),
			QuoteQuantity:   decimal.NewFromFloat(500),
			Commission:      decimal.NewFromFloat(0.5),
			CommissionAsset: "USDT",
			TradeTime:       tradeTime,
			IsBuyer:         true,
			IsMaker:         false,
		}

		mock.ExpectQuery("INSERT INTO trades").
			WithArgs(
				trade.UserID,
				trade.MexcTradeID,
				trade.OrderID,
				trade.Symbol,
				trade.Price,
				trade.Quantity,
				trade.QuoteQuantity,
				trade.Commission,
				trade.CommissionAsset,
				trade.TradeTime,
				trade.IsBuyer,
				trade.IsMaker,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, createdAt))

		err = s.CreateTrade(ctx, trade)
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), trade.ID)
		assert.Equal(t, createdAt, trade.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		trade := &domain.Trade{
			UserID:          1,
			MexcTradeID:     "123456",
			OrderID:         "order123",
			Symbol:          "BTCUSDT",
			Price:           decimal.NewFromFloat(50000),
			Quantity:        decimal.NewFromFloat(0.01),
			QuoteQuantity:   decimal.NewFromFloat(500),
			Commission:      decimal.NewFromFloat(0.5),
			CommissionAsset: "USDT",
			TradeTime:       time.Now(),
			IsBuyer:         true,
			IsMaker:         false,
		}

		mock.ExpectQuery("INSERT INTO trades").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		err = s.CreateTrade(ctx, trade)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTradeStorage_GetTradeByID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		mexcTradeID := "123456"
		tradeTime := time.Now().Add(-time.Hour)
		createdAt := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "mexc_trade_id", "order_id", "symbol", "price", "quantity",
			"quote_quantity", "commission", "commission_asset", "trade_time",
			"is_buyer", "is_maker", "created_at",
		}).AddRow(
			10, 1, mexcTradeID, "order123", "BTCUSDT", decimal.NewFromFloat(50000), decimal.NewFromFloat(0.01),
			decimal.NewFromFloat(500), decimal.NewFromFloat(0.5), "USDT", tradeTime,
			true, false, createdAt,
		)

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(mexcTradeID).
			WillReturnRows(rows)

		trade, err := s.GetTradeByID(ctx, mexcTradeID)
		assert.NoError(t, err)
		assert.NotNil(t, trade)
		assert.Equal(t, uint64(10), trade.ID)
		assert.Equal(t, mexcTradeID, trade.MexcTradeID)
		assert.Equal(t, "BTCUSDT", trade.Symbol)
		assert.Equal(t, "50000", trade.Price.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("trade not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		mexcTradeID := "123456"

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(mexcTradeID).
			WillReturnError(sql.ErrNoRows)

		trade, err := s.GetTradeByID(ctx, mexcTradeID)
		assert.Error(t, err)
		assert.Equal(t, postgreserr.ErrTradeNotFound, err)
		assert.Nil(t, trade)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		mexcTradeID := "123456"

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(mexcTradeID).
			WillReturnError(errors.New("db error"))

		trade, err := s.GetTradeByID(ctx, mexcTradeID)
		assert.Error(t, err)
		assert.Nil(t, trade)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTradeStorage_GetUserTrades(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get with filters", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		userID := uint64(1)
		startTime := time.Now().Add(-24 * time.Hour)
		endTime := time.Now()
		filter := TradeFilter{
			Symbol:    "BTCUSDT",
			StartTime: &startTime,
			EndTime:   &endTime,
			Limit:     10,
			Offset:    0,
		}

		tradeTime := time.Now().Add(-time.Hour)
		createdAt := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "mexc_trade_id", "order_id", "symbol", "price", "quantity",
			"quote_quantity", "commission", "commission_asset", "trade_time",
			"is_buyer", "is_maker", "created_at",
		}).AddRow(
			1, userID, "123456", "order123", "BTCUSDT", decimal.NewFromFloat(50000), decimal.NewFromFloat(0.01),
			decimal.NewFromFloat(500), decimal.NewFromFloat(0.5), "USDT", tradeTime,
			true, false, createdAt,
		)

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(userID, "BTCUSDT", startTime, endTime, 10).
			WillReturnRows(rows)

		trades, err := s.GetUserTrades(ctx, userID, filter)
		assert.NoError(t, err)
		assert.Len(t, trades, 1)
		assert.Equal(t, "BTCUSDT", trades[0].Symbol)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no filters", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		userID := uint64(1)

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "mexc_trade_id", "order_id", "symbol", "price", "quantity",
			"quote_quantity", "commission", "commission_asset", "trade_time",
			"is_buyer", "is_maker", "created_at",
		})

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(userID).
			WillReturnRows(rows)

		trades, err := s.GetUserTrades(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, trades, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		userID := uint64(1)

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(userID).
			WillReturnError(errors.New("query error"))

		trades, err := s.GetUserTrades(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, trades)
		assert.Contains(t, err.Error(), "query error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbAdapter := storage.NewDBAdapter(db)
		s := NewTradeStorage(dbAdapter)

		userID := uint64(1)

		// Return row with invalid data type
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "mexc_trade_id", "order_id", "symbol", "price", "quantity",
			"quote_quantity", "commission", "commission_asset", "trade_time",
			"is_buyer", "is_maker", "created_at",
		}).AddRow(
			"invalid_id", userID, "123456", "order123", "BTCUSDT", decimal.NewFromFloat(50000), decimal.NewFromFloat(0.01),
			decimal.NewFromFloat(500), decimal.NewFromFloat(0.5), "USDT", time.Now(),
			true, false, time.Now(),
		)

		mock.ExpectQuery("SELECT id, user_id, mexc_trade_id, order_id, symbol, price, quantity").
			WithArgs(userID).
			WillReturnRows(rows)

		_, err = s.GetUserTrades(ctx, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to scan trade")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
