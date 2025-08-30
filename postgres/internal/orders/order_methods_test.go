package orders

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres/postgreserr"
	"github.com/samar/sup_bot/metacore/storage/mocks"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// OrderStorageTestSuite содержит все тесты для OrderStorage
type OrderStorageTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockDB       *mocks.MockDBInterface
	orderStorage *OrderStorage
	ctx          context.Context
}

// SetupSuite вызывается один раз перед всеми тестами
func (suite *OrderStorageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest вызывается перед каждым тестом
func (suite *OrderStorageTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = mocks.NewMockDBInterface(suite.ctrl)
	suite.orderStorage = NewOrderStorage(suite.mockDB)
}

// TearDownTest вызывается после каждого теста
func (suite *OrderStorageTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestOrderStorage запускает все тесты
func TestOrderStorage(t *testing.T) {
	suite.Run(t, new(OrderStorageTestSuite))
}

// TestCreateOrder тестирует создание ордера
func (suite *OrderStorageTestSuite) TestCreateOrder() {
	order := &domain.Order{
		InternalID:          123,
		UserID:              1,
		MexcOrderID:         "mexc_order_123",
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.NewFromFloat(50000.00),
		Quantity:            decimal.NewFromFloat(0.001),
		QuoteOrderQty:       decimal.NewFromFloat(50.00),
		ExecutedQuantity:    decimal.NewFromFloat(0.000),
		CummulativeQuoteQty: decimal.NewFromFloat(0.00),
		ClientOrderID:       "client_123",
		TransactTime:        time.Now(),
	}

	suite.Run("successful creation", func() {
		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil)

		err := suite.orderStorage.CreateOrder(suite.ctx, order)

		assert.NoError(suite.T(), err)
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedError)

		err := suite.orderStorage.CreateOrder(suite.ctx, order)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to create order")
	})
}

// TestDeleteOrderByID тестирует удаление ордера
func (suite *OrderStorageTestSuite) TestDeleteOrderByID() {
	mexcOrderID := "mexc_order_123"

	suite.Run("successful deletion", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.DeleteOrderByID(suite.ctx, mexcOrderID)

		assert.NoError(suite.T(), err)
	})

	suite.Run("order not found", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(0), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.DeleteOrderByID(suite.ctx, mexcOrderID)

		assert.Error(suite.T(), err)
		assert.ErrorIs(suite.T(), err, postgreserr.ErrOrderNotFound)
		assert.Contains(suite.T(), err.Error(), "order with id mexc_order_123 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(nil, expectedError)

		err := suite.orderStorage.DeleteOrderByID(suite.ctx, mexcOrderID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to delete order")
	})

	suite.Run("rows affected error", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		expectedError := errors.New("failed to get rows affected")
		mockResult.EXPECT().RowsAffected().Return(int64(0), expectedError)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.DeleteOrderByID(suite.ctx, mexcOrderID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to get rows affected")
	})
}

// TestUpdateOrderStatus тестирует обновление статуса ордера
func (suite *OrderStorageTestSuite) TestUpdateOrderStatus() {
	mexcOrderID := "mexc_order_123"
	status := "FILLED"

	suite.Run("successful status update", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), status, mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.UpdateOrderStatus(suite.ctx, mexcOrderID, status)

		assert.NoError(suite.T(), err)
	})

	suite.Run("order not found", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(0), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), status, mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.UpdateOrderStatus(suite.ctx, mexcOrderID, status)

		assert.Error(suite.T(), err)
		assert.ErrorIs(suite.T(), err, postgreserr.ErrOrderNotFound)
		assert.Contains(suite.T(), err.Error(), "order with id mexc_order_123 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), status, mexcOrderID).
			Return(nil, expectedError)

		err := suite.orderStorage.UpdateOrderStatus(suite.ctx, mexcOrderID, status)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to update order status")
	})

	suite.Run("rows affected error", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		expectedError := errors.New("failed to get rows affected")
		mockResult.EXPECT().RowsAffected().Return(int64(0), expectedError)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), status, mexcOrderID).
			Return(mockResult, nil)

		err := suite.orderStorage.UpdateOrderStatus(suite.ctx, mexcOrderID, status)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to get rows affected")
	})
}

// TestGetOrderByID тестирует получение ордера по ID
func (suite *OrderStorageTestSuite) TestGetOrderByID() {
	mexcOrderID := "mexc_order_123"

	suite.Run("successful retrieval", func() {
		// Создаем мок для RowInterface
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(dest ...interface{}) error {
				// Устанавливаем значения в dest
				if len(dest) >= 1 {
					if id, ok := dest[0].(*uint64); ok {
						*id = 1
					}
				}
				if len(dest) >= 2 {
					if internalID, ok := dest[1].(*int64); ok {
						*internalID = 123
					}
				}
				if len(dest) >= 3 {
					if userID, ok := dest[2].(*uint64); ok {
						*userID = 1
					}
				}
				if len(dest) >= 4 {
					if mexcOrderID, ok := dest[3].(*string); ok {
						*mexcOrderID = "mexc_order_123"
					}
				}
				if len(dest) >= 5 {
					if symbol, ok := dest[4].(*string); ok {
						*symbol = "BTCUSDT"
					}
				}
				if len(dest) >= 6 {
					if side, ok := dest[5].(*string); ok {
						*side = "BUY"
					}
				}
				// Устанавливаем остальные поля...
				return nil
			})

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockRow)

		order, err := suite.orderStorage.GetOrderByID(suite.ctx, mexcOrderID)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), order)
		assert.Equal(suite.T(), mexcOrderID, order.MexcOrderID)
		assert.Equal(suite.T(), "BTCUSDT", order.Symbol)
		assert.Equal(suite.T(), "BUY", order.Side)
	})

	suite.Run("order not found", func() {
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockRow)

		order, err := suite.orderStorage.GetOrderByID(suite.ctx, mexcOrderID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), order)
		assert.Contains(suite.T(), err.Error(), "order with id mexc_order_123 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedError)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), mexcOrderID).
			Return(mockRow)

		order, err := suite.orderStorage.GetOrderByID(suite.ctx, mexcOrderID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), order)
		assert.Contains(suite.T(), err.Error(), "failed to get order")
	})
}

// TestNewOrderStorage тестирует создание нового экземпляра OrderStorage
func (suite *OrderStorageTestSuite) TestNewOrderStorage() {
	assert.NotNil(suite.T(), suite.orderStorage)
	assert.Equal(suite.T(), suite.mockDB, suite.orderStorage.db)
}

// TestClose тестирует закрытие соединения
func (suite *OrderStorageTestSuite) TestClose() {
	suite.mockDB.EXPECT().Close().Return(nil)
	suite.orderStorage.Close()
}
