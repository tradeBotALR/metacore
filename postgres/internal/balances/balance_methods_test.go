package balances

import (
	"context"
	"testing"
	"time"

	"metacore/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BalanceStorageTestSuite содержит все тесты для BalanceStorage
type BalanceStorageTestSuite struct {
	suite.Suite
	balanceStorage *BalanceStorage
	ctx            context.Context
}

// SetupSuite вызывается один раз перед всеми тестами
func (suite *BalanceStorageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest вызывается перед каждым тестом
func (suite *BalanceStorageTestSuite) SetupTest() {
	suite.balanceStorage = NewBalanceStorage()
}

// TestBalanceStorage запускает все тесты
func TestBalanceStorage(t *testing.T) {
	suite.Run(t, new(BalanceStorageTestSuite))
}

// TestUpdateBalance тестирует обновление баланса
func (suite *BalanceStorageTestSuite) TestUpdateBalance() {
	balance := &domain.UserBalance{
		UserID:    1,
		Asset:     "BTC",
		Free:      decimal.NewFromFloat(0.4),
		Locked:    decimal.NewFromFloat(0.1),
		UpdatedAt: time.Now(),
	}

	suite.Run("successful update", func() {
		// Поскольку методы еще не реализованы, тестируем только структуру
		err := suite.balanceStorage.UpdateBalance(suite.ctx, balance)

		assert.NoError(suite.T(), err)
	})
}

// TestGetBalance тестирует получение баланса
func (suite *BalanceStorageTestSuite) TestGetBalance() {
	userID := uint64(1)
	asset := "BTC"

	suite.Run("successful retrieval", func() {
		// Поскольку методы еще не реализованы, тестируем только структуру
		balance, err := suite.balanceStorage.GetBalance(suite.ctx, userID, asset)

		assert.NoError(suite.T(), err)
		assert.Nil(suite.T(), balance) // Пока возвращается nil, так как метод не реализован
	})
}

// TestNewBalanceStorage тестирует создание нового экземпляра BalanceStorage
func (suite *BalanceStorageTestSuite) TestNewBalanceStorage() {
	assert.NotNil(suite.T(), suite.balanceStorage)
}
