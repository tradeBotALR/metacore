package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"

	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres"
)

// RunOrderScenario демонстрирует базовый сценарий работы с библиотекой:
// 1) Создание пользователя
// 2) Создание ордера
// 3) Получение ордера
// 4) Обновление статуса ордера (FILLED)
// 5) Создание сделки по ордеру
// 6) Обновление балансов пользователя
// 7) Удаление ордера
func RunOrderScenario(ctx context.Context, db *postgres.DB) error {
	log.Println("\n🚦 Запуск order-сценария...")

	// Шаг 1. Создание пользователя
	user := &domain.User{
		TelegramID:    time.Now().Unix(),
		MexcUID:       fmt.Sprintf("demo_uid_%d", time.Now().UnixNano()),
		Username:      "scenario_user",
		Email:         "scenario_user@example.com",
		MexcAPIKey:    "api_key_placeholder",
		MexcSecretKey: "secret_key_placeholder",
		KYCStatus:     1,
		CanTrade:      true,
		CanWithdraw:   true,
		CanDeposit:    true,
		AccountType:   "spot",
		Permissions:   `{"perms":["trade"]}`,
		IsActive:      true,
	}
	if err := db.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	log.Printf("👤 Пользователь создан: id=%d, username=%s", user.ID, user.Username)

	// Шаг 2. Создание ордера
	mexcOrderID := fmt.Sprintf("order_%d", time.Now().UnixNano())
	order := &domain.Order{
		InternalID:          time.Now().UnixNano(),
		UserID:              user.ID,
		MexcOrderID:         mexcOrderID,
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.NewFromFloat(50000.00),
		Quantity:            decimal.NewFromFloat(0.002),
		QuoteOrderQty:       decimal.NewFromFloat(100.00),
		ExecutedQuantity:    decimal.Zero,
		CummulativeQuoteQty: decimal.Zero,
		ClientOrderID:       fmt.Sprintf("cli_%d", time.Now().UnixNano()),
		TransactTime:        time.Now(),
	}
	if err := db.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	log.Printf("📈 Ордер создан: mexc_order_id=%s, symbol=%s, side=%s", order.MexcOrderID, order.Symbol, order.Side)

	// Шаг 3. Получение ордера по mexc_order_id
	gotOrder, err := db.GetOrderByID(ctx, mexcOrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}
	log.Printf("🔎 Ордер найден: %s %s %s qty=%s", gotOrder.Symbol, gotOrder.Side, gotOrder.Status, gotOrder.Quantity.String())

	// Шаг 4. Обновление статуса ордера -> FILLED
	if err := db.UpdateOrderStatus(ctx, mexcOrderID, "FILLED"); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	log.Println("🔄 Статус ордера обновлен на FILLED")

	// Шаг 5. Создание сделки по ордеру
	trade := &domain.Trade{
		UserID:          user.ID,
		MexcTradeID:     fmt.Sprintf("trade_%d", time.Now().UnixNano()),
		OrderID:         mexcOrderID,
		Symbol:          "BTCUSDT",
		Price:           decimal.NewFromFloat(50000.00),
		Quantity:        decimal.NewFromFloat(0.002),
		QuoteQuantity:   decimal.NewFromFloat(100.00),
		Commission:      decimal.NewFromFloat(0.01),
		CommissionAsset: "USDT",
		TradeTime:       time.Now(),
		IsBuyer:         true,
		IsMaker:         false,
	}
	if err := db.CreateTrade(ctx, trade); err != nil {
		return fmt.Errorf("create trade: %w", err)
	}
	log.Printf("🤝 Сделка создана: mexc_trade_id=%s qty=%s", trade.MexcTradeID, trade.Quantity.String())

	// Шаг 6. Обновление балансов пользователя (пример: списали USDT, зачислили BTC)
	balances := []*domain.UserBalance{
		{UserID: user.ID, Asset: "USDT", Free: decimal.NewFromFloat(900.00), Locked: decimal.Zero},
		{UserID: user.ID, Asset: "BTC", Free: decimal.NewFromFloat(0.002), Locked: decimal.Zero},
	}
	if err := db.UpdateUserBalances(ctx, user.ID, balances); err != nil {
		return fmt.Errorf("update user balances: %w", err)
	}
	log.Println("💰 Балансы пользователя обновлены")

	// Шаг 7. Удаление ордера (как завершенного)
	if err := db.DeleteOrderByID(ctx, mexcOrderID); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	log.Println("🗑️  Ордер удален (завершен)")

	log.Println("✅ Order-сценарий успешно выполнен")
	return nil
}
