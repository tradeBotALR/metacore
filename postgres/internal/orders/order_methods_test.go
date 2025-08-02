package orders

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"metacore/domain"
	"metacore/storage"
	"metacore/storage/mocks"
	"testing"
	"time"
)

var store storage.FullStorage

func TestStorage_CreateOrder(t *testing.T) {
	// 1. Создаем контроллер для gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Убедимся, что все ожидания будут проверены

	// 2. Создаем мок пула
	mockPool := mocks.NewMockPgxPoolIface(ctrl)

	// 3. Создаем экземпляр нашего хранилища с моком
	s := NewOrderStorage(mockPool)

	// 4. Подготавливаем тестовые данные
	ctx := context.Background()
	testOrder := &domain.Order{
		InternalID:          1001,
		UserID:              123,
		MexcOrderID:         "test_mexc_order_id_123",
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.RequireFromString("50000.0"),
		Quantity:            decimal.RequireFromString("0.001"),
		QuoteOrderQty:       decimal.Zero,
		ExecutedQuantity:    decimal.Zero,
		CummulativeQuoteQty: decimal.Zero,
		ClientOrderID:       "my_client_order_1",
		TransactTime:        time.Now().UTC(), // Используем time.Time
	}

	// --- 2. Тест: Создание ордера (CreateOrder) ---
	mockPool.EXPECT().
		Exec(ctx, // SQL строка проверяется через gomock.Any() для простоты
			gomock.Any(), // SQL
			testOrder.InternalID,
			testOrder.UserID,
			testOrder.MexcOrderID,
			testOrder.Symbol,
			testOrder.Side,
			testOrder.Type,
			testOrder.Status,
			testOrder.Price,
			testOrder.Quantity,
			testOrder.QuoteOrderQty,
			testOrder.ExecutedQuantity,
			testOrder.CummulativeQuoteQty,
			testOrder.ClientOrderID,
			testOrder.TransactTime,
		).
		DoAndReturn(func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
			// Можно добавить более точную проверку SQL, если нужно
			assert.Contains(t, sql, "INSERT INTO orders")
			return pgconn.CommandTag{}, nil // Возвращаем успех
		})

	err := s.CreateOrder(ctx, testOrder)
	// Проверяем результат
	assert.NoError(t, err, "CreateOrder should not return an error")

}

func TestOrderStorage_SimpleFlow_PGXMock(t *testing.T) {
	// --- 1. Настройка ---
	// Создаем мок пула с помощью pgxmock
	// pgxmock.PoolIface реализует интерфейс, похожий на pgxpool.Pool
	mockPool, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close() // Закрываем мок-соединение

	// Создаем экземпляр тестируемой реализации
	orderStorage := NewOrderStorage(mockPool) // mockPool реализует storage.PgxPoolIface

	// Тестовые данные
	ctx := context.Background()
	testOrder := &domain.Order{
		InternalID:          123,
		UserID:              1,
		MexcOrderID:         "test_mexc_id_12345",
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.RequireFromString("60000.0"),
		Quantity:            decimal.RequireFromString("0.001"),
		QuoteOrderQty:       decimal.Zero,
		ExecutedQuantity:    decimal.Zero,
		CummulativeQuoteQty: decimal.Zero,
		ClientOrderID:       "my_client_order_1",
		TransactTime:        time.Now().Truncate(time.Millisecond).UTC(),
	}

	// --- 2. Тест: Создание ордера (CreateOrder) ---
	// Ожидаем, что будет выполнен один INSERT
	mockPool.ExpectExec("INSERT INTO orders").
		WithArgs(
			// Порядок аргументов должен строго соответствовать запросу в CreateOrder
			testOrder.InternalID,
			testOrder.UserID,
			testOrder.MexcOrderID,
			testOrder.Symbol,
			testOrder.Side,
			testOrder.Type,
			testOrder.Status,
			testOrder.Price,
			testOrder.Quantity,
			testOrder.QuoteOrderQty,
			testOrder.ExecutedQuantity,
			testOrder.CummulativeQuoteQty,
			testOrder.ClientOrderID,
			testOrder.TransactTime,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)) // 1 row affected

	// Выполняем тестируемую функцию
	err = orderStorage.CreateOrder(ctx, testOrder)
	// Проверяем результат
	assert.NoError(t, err)

	// Проверяем, что все ожидания были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())

	// --- 3. Тест: Получение ордера (GetOrderByID) ---
	// Ожидаем, что будет выполнен SELECT
	rows := pgxmock.NewRows([]string{
		"id", "internal_id", "user_id", "mexc_order_id", "symbol", "side", "type", "status",
		"price", "quantity", "quote_order_qty", "executed_quantity",
		"cummulative_quote_qty", "client_order_id", "transact_time",
		"created_at", "updated_at",
	}).AddRow(
		uint64(1), testOrder.InternalID, testOrder.UserID, testOrder.MexcOrderID, testOrder.Symbol,
		testOrder.Side, testOrder.Type, testOrder.Status,
		testOrder.Price, testOrder.Quantity, testOrder.QuoteOrderQty,
		testOrder.ExecutedQuantity, testOrder.CummulativeQuoteQty,
		testOrder.ClientOrderID, testOrder.TransactTime,
		time.Now(), time.Now(), // created_at, updated_at
	)

	mockPool.ExpectQuery("SELECT (.+) FROM orders WHERE mexc_order_id").
		WithArgs(testOrder.MexcOrderID).
		WillReturnRows(rows)

	// Выполняем тестируемую функцию
	retrievedOrder, err := orderStorage.GetOrderByID(ctx, testOrder.MexcOrderID)
	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	if retrievedOrder != nil {
		assert.Equal(t, testOrder.MexcOrderID, retrievedOrder.MexcOrderID)
		assert.Equal(t, testOrder.Symbol, retrievedOrder.Symbol)
		// ... другие проверки ...
	}

	// Проверяем, что все ожидания были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())

	// --- 4. Тест: Удаление ордера (DeleteOrderByID) ---
	// Ожидаем, что будет выполнен DELETE
	mockPool.ExpectExec("DELETE FROM orders WHERE mexc_order_id").
		WithArgs(testOrder.MexcOrderID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 row affected

	// Выполняем тестируемую функцию
	err = orderStorage.DeleteOrderByID(ctx, testOrder.MexcOrderID)
	// Проверяем результат
	assert.NoError(t, err)

	// Проверяем, что все ожидания были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

// TestOrderStorage_SimpleFlow тестирует простой сценарий: создать -> получить -> удалить.
func TestOrderStorage_SimpleFlow(t *testing.T) {
	// --- 1. Настройка ---
	// Создаем мок пула с помощью pgxmock
	mockPool, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close()

	// Создаем экземпляр тестируемой реализации
	orderStorage := NewOrderStorage(mockPool)

	// Тестовые данные
	ctx := context.Background()
	testOrder := &domain.Order{
		InternalID:          555,
		UserID:              1,
		MexcOrderID:         "test_mexc_id_12345",
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.RequireFromString("60000.0"),
		Quantity:            decimal.RequireFromString("0.001"),
		QuoteOrderQty:       decimal.Zero,
		ExecutedQuantity:    decimal.Zero,
		CummulativeQuoteQty: decimal.Zero,
		ClientOrderID:       "my_client_order_1",
		TransactTime:        time.Now().Truncate(time.Millisecond).UTC(),
	}

	// --- 2. Тест: Создание ордера (CreateOrder) ---
	// Ожидаем, что будет выполнен один INSERT с определенными аргументами
	mockPool.ExpectExec("INSERT INTO orders").
		WithArgs(
			// Порядок аргументов должен строго соответствовать запросу в CreateOrder
			testOrder.InternalID,
			testOrder.UserID,
			testOrder.MexcOrderID,
			testOrder.Symbol,
			testOrder.Side,
			testOrder.Type,
			testOrder.Status,
			testOrder.Price,
			testOrder.Quantity,
			testOrder.QuoteOrderQty,
			testOrder.ExecutedQuantity,
			testOrder.CummulativeQuoteQty,
			testOrder.ClientOrderID,
			testOrder.TransactTime,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1)) // 1 row affected

	// Выполняем тестируемую функцию
	err = orderStorage.CreateOrder(ctx, testOrder)

	// Проверяем результат
	assert.NoError(t, err)

	// Убедимся, что все ожидания от мока были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())

	// --- 3. Тест: Получение ордера (GetOrderByID) ---
	// Подготавливаем "фиктивный" результат для SELECT запроса
	rows := pgxmock.NewRows([]string{
		"id", "internal_id", "user_id", "mexc_order_id", "symbol", "side", "type", "status",
		"price", "quantity", "quote_order_qty", "executed_quantity",
		"cummulative_quote_qty", "client_order_id", "transact_time",
		"created_at", "updated_at",
	}).AddRow(
		uint64(1), testOrder.InternalID, testOrder.UserID, testOrder.MexcOrderID, testOrder.Symbol,
		testOrder.Side, testOrder.Type, testOrder.Status,
		testOrder.Price, testOrder.Quantity, testOrder.QuoteOrderQty,
		testOrder.ExecutedQuantity, testOrder.CummulativeQuoteQty,
		testOrder.ClientOrderID, testOrder.TransactTime,
		time.Now(), time.Now(), // created_at, updated_at
	)

	// Ожидаем, что будет выполнен SELECT с определенным ID
	mockPool.ExpectQuery("SELECT (.+) FROM orders WHERE mexc_order_id").
		WithArgs(testOrder.MexcOrderID).
		WillReturnRows(rows)

	// Выполняем тестируемую функцию
	retrievedOrder, err := orderStorage.GetOrderByID(ctx, testOrder.MexcOrderID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	// Проверяем ключевые поля
	if retrievedOrder != nil {
		assert.Equal(t, testOrder.MexcOrderID, retrievedOrder.MexcOrderID)
		assert.Equal(t, testOrder.Symbol, retrievedOrder.Symbol)
		// ... можно добавить больше проверок assert.Equal для других полей ...
	}

	// Убедимся, что все ожидания от мока были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())

	// --- 4. Тест: Удаление ордера (DeleteOrderByID) ---
	// Ожидаем, что будет выполнен DELETE с определенным ID
	mockPool.ExpectExec("DELETE FROM orders WHERE mexc_order_id").
		WithArgs(testOrder.MexcOrderID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 row affected

	// Выполняем тестируемую функцию
	err = orderStorage.DeleteOrderByID(ctx, testOrder.MexcOrderID)

	// Проверяем результат
	assert.NoError(t, err)

	// Убедимся, что все ожидания от мока были выполнены
	assert.NoError(t, mockPool.ExpectationsWereMet())
}
