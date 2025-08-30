package users

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserStorageTestSuite содержит все тесты для UserStorage
type UserStorageTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockDB      *mocks.MockDBInterface
	userStorage *UserStorage
	ctx         context.Context
}

// SetupSuite вызывается один раз перед всеми тестами
func (suite *UserStorageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest вызывается перед каждым тестом
func (suite *UserStorageTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = mocks.NewMockDBInterface(suite.ctrl)
	suite.userStorage = NewUserStorage(suite.mockDB)
}

// TearDownTest вызывается после каждого теста
func (suite *UserStorageTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestUserStorage запускает все тесты
func TestUserStorage(t *testing.T) {
	suite.Run(t, new(UserStorageTestSuite))
}

// TestCreateUser тестирует создание пользователя
func (suite *UserStorageTestSuite) TestCreateUser() {
	user := &domain.User{
		MexcUID:       "test_uid_123",
		Username:      "testuser",
		Email:         "test@example.com",
		MexcAPIKey:    "api_key_123",
		MexcSecretKey: "secret_key_123",
		KYCStatus:     1, // verified
		CanTrade:      true,
		CanWithdraw:   true,
		CanDeposit:    true,
		AccountType:   "spot",
		Permissions:   `["trade", "withdraw"]`, // JSON string
		IsActive:      true,
	}

	suite.Run("successful creation", func() {
		// Создаем мок для RowInterface
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(dest ...interface{}) error {
				// Устанавливаем значения в dest
				if len(dest) >= 1 {
					if id, ok := dest[0].(*uint64); ok {
						*id = 1
					}
				}
				if len(dest) >= 2 {
					if createdAt, ok := dest[1].(*time.Time); ok {
						*createdAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
					}
				}
				if len(dest) >= 3 {
					if updatedAt, ok := dest[2].(*time.Time); ok {
						*updatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
					}
				}
				return nil
			})

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockRow)

		err := suite.userStorage.CreateUser(suite.ctx, user)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), uint64(1), user.ID)
		expectedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(suite.T(), expectedTime, user.CreatedAt)
		assert.Equal(suite.T(), expectedTime, user.UpdatedAt)
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedError)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockRow)

		err := suite.userStorage.CreateUser(suite.ctx, user)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to create user")
	})
}

// TestGetUserByID тестирует получение пользователя по ID
func (suite *UserStorageTestSuite) TestGetUserByID() {
	userID := uint64(1)

	suite.Run("successful retrieval", func() {
		// Создаем мок для RowInterface
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(dest ...interface{}) error {
				// Устанавливаем значения в dest
				if len(dest) >= 1 {
					if id, ok := dest[0].(*uint64); ok {
						*id = 1
					}
				}
				if len(dest) >= 2 {
					if mexcUID, ok := dest[1].(*string); ok {
						*mexcUID = "test_uid_123"
					}
				}
				if len(dest) >= 3 {
					if username, ok := dest[2].(*string); ok {
						*username = "testuser"
					}
				}
				if len(dest) >= 4 {
					if email, ok := dest[3].(*string); ok {
						*email = "test@example.com"
					}
				}
				// Устанавливаем остальные поля...
				return nil
			})

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), userID).
			Return(mockRow)

		user, err := suite.userStorage.GetUserByID(suite.ctx, userID)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), user)
		assert.Equal(suite.T(), userID, user.ID)
		assert.Equal(suite.T(), "testuser", user.Username)
		assert.Equal(suite.T(), "test@example.com", user.Email)
	})

	suite.Run("user not found", func() {
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), userID).
			Return(mockRow)

		user, err := suite.userStorage.GetUserByID(suite.ctx, userID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
		assert.Contains(suite.T(), err.Error(), "user with id 1 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedError)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), userID).
			Return(mockRow)

		user, err := suite.userStorage.GetUserByID(suite.ctx, userID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
		assert.Contains(suite.T(), err.Error(), "failed to get user")
	})
}

// TestGetUserByMexcUID тестирует получение пользователя по MEXC UID
func (suite *UserStorageTestSuite) TestGetUserByMexcUID() {
	mexcUID := "test_uid_123"

	suite.Run("successful retrieval", func() {
		// Создаем мок для RowInterface
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(dest ...interface{}) error {
				// Устанавливаем значения в dest
				if len(dest) >= 1 {
					if id, ok := dest[0].(*uint64); ok {
						*id = 1
					}
				}
				if len(dest) >= 2 {
					if mexcUID, ok := dest[1].(*string); ok {
						*mexcUID = "test_uid_123"
					}
				}
				if len(dest) >= 3 {
					if username, ok := dest[2].(*string); ok {
						*username = "testuser"
					}
				}
				// Устанавливаем остальные поля...
				return nil
			})

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), mexcUID).
			Return(mockRow)

		user, err := suite.userStorage.GetUserByMexcUID(suite.ctx, mexcUID)

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), user)
		assert.Equal(suite.T(), mexcUID, user.MexcUID)
		assert.Equal(suite.T(), "testuser", user.Username)
	})

	suite.Run("user not found", func() {
		mockRow := mocks.NewMockRowInterface(suite.ctrl)
		mockRow.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		suite.mockDB.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any(), mexcUID).
			Return(mockRow)

		user, err := suite.userStorage.GetUserByMexcUID(suite.ctx, mexcUID)

		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
		assert.Contains(suite.T(), err.Error(), "user with mexc_uid test_uid_123 not found")
	})
}

// TestUpdateUser тестирует обновление пользователя
func (suite *UserStorageTestSuite) TestUpdateUser() {
	user := &domain.User{
		ID:            1,
		MexcUID:       "test_uid_123",
		Username:      "updated_user",
		Email:         "updated@example.com",
		MexcAPIKey:    "new_api_key",
		MexcSecretKey: "new_secret_key",
		KYCStatus:     0, // pending
		CanTrade:      false,
		CanWithdraw:   false,
		CanDeposit:    true,
		AccountType:   "futures",
		Permissions:   `["deposit"]`, // JSON string
		IsActive:      false,
	}

	suite.Run("successful update", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockResult, nil)

		err := suite.userStorage.UpdateUser(suite.ctx, user)

		assert.NoError(suite.T(), err)
	})

	suite.Run("user not found", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(0), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockResult, nil)

		err := suite.userStorage.UpdateUser(suite.ctx, user)

		assert.Error(suite.T(), err)
		assert.ErrorIs(suite.T(), err, postgreserr.ErrUserNotFound)
		assert.Contains(suite.T(), err.Error(), "user with id 1 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedError)

		err := suite.userStorage.UpdateUser(suite.ctx, user)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to update user")
	})

	suite.Run("rows affected error", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		expectedError := errors.New("failed to get rows affected")
		mockResult.EXPECT().RowsAffected().Return(int64(0), expectedError)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockResult, nil)

		err := suite.userStorage.UpdateUser(suite.ctx, user)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to get rows affected")
	})
}

// TestDeleteUser тестирует удаление пользователя
func (suite *UserStorageTestSuite) TestDeleteUser() {
	userID := uint64(1)

	suite.Run("successful deletion", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), userID).
			Return(mockResult, nil)

		err := suite.userStorage.DeleteUser(suite.ctx, userID)

		assert.NoError(suite.T(), err)
	})

	suite.Run("user not found", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		mockResult.EXPECT().RowsAffected().Return(int64(0), nil)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), userID).
			Return(mockResult, nil)

		err := suite.userStorage.DeleteUser(suite.ctx, userID)

		assert.Error(suite.T(), err)
		assert.ErrorIs(suite.T(), err, postgreserr.ErrUserNotFound)
		assert.Contains(suite.T(), err.Error(), "user with id 1 not found")
	})

	suite.Run("database error", func() {
		expectedError := errors.New("database connection failed")

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), userID).
			Return(nil, expectedError)

		err := suite.userStorage.DeleteUser(suite.ctx, userID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to delete user")
	})

	suite.Run("rows affected error", func() {
		mockResult := mocks.NewMockResult(suite.ctrl)
		expectedError := errors.New("failed to get rows affected")
		mockResult.EXPECT().RowsAffected().Return(int64(0), expectedError)

		suite.mockDB.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), userID).
			Return(mockResult, nil)

		err := suite.userStorage.DeleteUser(suite.ctx, userID)

		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "failed to get rows affected")
	})
}

// TestNewUserStorage тестирует создание нового экземпляра UserStorage
func (suite *UserStorageTestSuite) TestNewUserStorage() {
	assert.NotNil(suite.T(), suite.userStorage)
	assert.Equal(suite.T(), suite.mockDB, suite.userStorage.db)
}

// TestClose тестирует закрытие соединения
func (suite *UserStorageTestSuite) TestClose() {
	suite.mockDB.EXPECT().Close().Return(nil)
	suite.userStorage.Close()
}
