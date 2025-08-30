# Metacore Library

`metacore` - это Go библиотека для работы с базой данных, пользователями, ордерами, сделками и балансами.

## Установка

```bash
go get github.com/samar/sup_bot/metacore
```

## Использование

### Импорт

```go
import "github.com/samar/sup_bot/metacore"
```

### Подключение к базе данных

```go
// Получаем конфигурацию по умолчанию
cfg := metacore.DefaultConfig()

// Создаем подключение к PostgreSQL
db, err := metacore.NewPostgresDB(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Работа с пользователями

```go
// Создание пользователя
user := &metacore.User{
    MexcUID:       "user_123",
    Username:      "trader",
    Email:         "trader@example.com",
    MexcAPIKey:    "api_key",
    MexcSecretKey: "secret_key",
    KYCStatus:     1,
    CanTrade:      true,
    CanWithdraw:   true,
    CanDeposit:    true,
    AccountType:   "spot",
    Permissions:   `["trade", "withdraw", "deposit"]`,
    IsActive:      true,
}

err := db.CreateUser(ctx, user)
if err != nil {
    log.Printf("Error creating user: %v", err)
}

// Получение пользователя по ID
retrievedUser, err := db.GetUserByID(ctx, user.ID)
if err != nil {
    log.Printf("Error getting user: %v", err)
}

// Получение пользователя по MEXC UID
userByUID, err := db.GetUserByMexcUID(ctx, "user_123")
if err != nil {
    log.Printf("Error getting user by UID: %v", err)
}
```

### Работа с ордерами

```go
// Создание ордера
order := &metacore.Order{
    InternalID:          12345,
    UserID:              user.ID,
    MexcOrderID:         "order_123",
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

err := db.CreateOrder(ctx, order)
if err != nil {
    log.Printf("Error creating order: %v", err)
}

// Получение ордера по ID
retrievedOrder, err := db.GetOrderByID(ctx, "order_123")
if err != nil {
    log.Printf("Error getting order: %v", err)
}

// Обновление статуса ордера
err = db.UpdateOrderStatus(ctx, "order_123", "FILLED")
if err != nil {
    log.Printf("Error updating order status: %v", err)
}
```

### Работа со сделками

```go
// Создание сделки
trade := &metacore.Trade{
    InternalID:    12345,
    UserID:        user.ID,
    MexcTradeID:   "trade_123",
    OrderID:       "order_123",
    Symbol:        "BTCUSDT",
    Side:          "BUY",
    Price:         decimal.NewFromFloat(50000.00),
    Quantity:      decimal.NewFromFloat(0.001),
    QuoteQuantity: decimal.NewFromFloat(50.00),
    Commission:    decimal.NewFromFloat(0.05),
    CommissionAsset: "USDT",
    TransactTime:  time.Now(),
}

err := db.CreateTrade(ctx, trade)
if err != nil {
    log.Printf("Error creating trade: %v", err)
}

// Получение сделки по ID
retrievedTrade, err := db.GetTradeByID(ctx, "trade_123")
if err != nil {
    log.Printf("Error getting trade: %v", err)
}
```

### Работа с балансами

```go
// Обновление баланса
balance := &metacore.UserBalance{
    UserID:    user.ID,
    Asset:     "BTC",
    Balance:   decimal.NewFromFloat(0.5),
    Available: decimal.NewFromFloat(0.4),
    Frozen:    decimal.NewFromFloat(0.1),
}

err := db.UpdateBalance(ctx, balance)
if err != nil {
    log.Printf("Error updating balance: %v", err)
}

// Получение баланса
retrievedBalance, err := db.GetBalance(ctx, user.ID, "BTC")
if err != nil {
    log.Printf("Error getting balance: %v", err)
}
```

## Интерфейсы

Библиотека предоставляет следующие интерфейсы:

- `FullStorage` - объединяет все интерфейсы хранилища
- `UserStorage` - для работы с пользователями
- `OrderStorage` - для работы с ордерами
- `TradeStorage` - для работы со сделками
- `BalanceStorage` - для работы с балансами
- `DBInterface` - для работы с базой данных

## Конфигурация

По умолчанию используется следующая конфигурация:

```go
type Config struct {
    DB struct {
        Host     string
        Port     int
        User     string
        Password string
        DBName   string
        SSLMode  string
    }
    Pool struct {
        MaxConns         int
        MinConns         int
        MaxConnLifetime  time.Duration
        MaxConnIdleTime  time.Duration
    }
}
```

## Зависимости

- `github.com/lib/pq` - PostgreSQL драйвер
- `github.com/shopspring/decimal` - для работы с десятичными числами
- `github.com/golang/mock` - для создания моков в тестах
- `github.com/stretchr/testify` - для тестирования

## Тестирование

```bash
cd metacore
go test ./...
```

## Примеры

Полные примеры использования можно найти в файле `cmd/main.go`.
