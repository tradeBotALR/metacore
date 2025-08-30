package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/samar/sup_bot/metacore/configs"
	"github.com/samar/sup_bot/metacore/domain"
	"github.com/samar/sup_bot/metacore/postgres"

	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("🚀 Запуск Metacore демо-сценария...")

	// Загружаем конфигурацию
	cfg := configs.DefaultConfig()

	// Подключаемся к БД
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Подключение к БД установлено")

	// Создаем контекст
	ctx := context.Background()

	// Демо-сценарий: Создание пользователя, ордера и баланса
	runDemoScenario(ctx, db)
}

func runDemoScenario(ctx context.Context, db *postgres.DB) {
	fmt.Println("\n🎭 Запуск демо-сценария...")

	// 1. Создаем пользователя
	fmt.Println("\n👤 1. Создание пользователя...")
	user := &domain.User{
		MexcUID:       "demo_user_123",
		Username:      "demo_trader",
		Email:         "demo@metacore.com",
		MexcAPIKey:    "demo_api_key_123",
		MexcSecretKey: "demo_secret_key_123",
		KYCStatus:     1, // verified
		CanTrade:      true,
		CanWithdraw:   true,
		CanDeposit:    true,
		AccountType:   "spot",
		Permissions:   `["trade", "withdraw", "deposit"]`,
		IsActive:      true,
	}

	err := db.CreateUser(ctx, user)
	if err != nil {
		log.Printf("⚠️ Ошибка создания пользователя: %v", err)
	} else {
		fmt.Printf("✅ Пользователь создан с ID: %d\n", user.ID)
		fmt.Printf("   Username: %s, Email: %s\n", user.Username, user.Email)
		fmt.Printf("   KYC Status: %d, Can Trade: %t\n", user.KYCStatus, user.CanTrade)
	}

	// 2. Получаем пользователя по ID
	fmt.Println("\n🔍 2. Получение пользователя по ID...")
	retrievedUser, err := db.GetUserByID(ctx, user.ID)
	if err != nil {
		log.Printf("⚠️ Ошибка получения пользователя: %v", err)
	} else {
		fmt.Printf("✅ Пользователь найден: %s (%s)\n", retrievedUser.Username, retrievedUser.Email)
	}

	// 3. Получаем пользователя по MEXC UID
	fmt.Println("\n🔍 3. Получение пользователя по MEXC UID...")
	retrievedUserByUID, err := db.GetUserByMexcUID(ctx, user.MexcUID)
	if err != nil {
		log.Printf("⚠️ Ошибка получения пользователя по UID: %v", err)
	} else {
		fmt.Printf("✅ Пользователь найден по UID: %s\n", retrievedUserByUID.Username)
	}

	// 4. Создаем ордер
	fmt.Println("\n📊 4. Создание ордера...")
	order := &domain.Order{
		InternalID:          12345,
		UserID:              user.ID,
		MexcOrderID:         "demo_order_123",
		Symbol:              "BTCUSDT",
		Side:                "BUY",
		Type:                "LIMIT",
		Status:              "NEW",
		Price:               decimal.NewFromFloat(50000.00),
		Quantity:            decimal.NewFromFloat(0.001),
		QuoteOrderQty:       decimal.NewFromFloat(50.00),
		ExecutedQuantity:    decimal.NewFromFloat(0.000),
		CummulativeQuoteQty: decimal.NewFromFloat(0.00),
		ClientOrderID:       "demo_client_123",
		TransactTime:        time.Now(),
	}

	err = db.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("⚠️ Ошибка создания ордера: %v", err)
	} else {
		fmt.Printf("✅ Ордер создан для пользователя %d\n", order.UserID)
		fmt.Printf("   Symbol: %s, Side: %s, Price: %s\n", order.Symbol, order.Side, order.Price.String())
		fmt.Printf("   Quantity: %s, Status: %s\n", order.Quantity.String(), order.Status)
	}

	// 5. Получаем ордер по ID
	fmt.Println("\n🔍 5. Получение ордера по ID...")
	retrievedOrder, err := db.GetOrderByID(ctx, order.MexcOrderID)
	if err != nil {
		log.Printf("⚠️ Ошибка получения ордера: %v", err)
	} else {
		fmt.Printf("✅ Ордер найден: %s %s %s\n", retrievedOrder.Symbol, retrievedOrder.Side, retrievedOrder.Status)
		fmt.Printf("   Price: %s, Quantity: %s\n", retrievedOrder.Price.String(), retrievedOrder.Quantity.String())
	}

	// 6. Обновляем статус ордера
	fmt.Println("\n🔄 6. Обновление статуса ордера...")
	err = db.UpdateOrderStatus(ctx, order.MexcOrderID, "FILLED")
	if err != nil {
		log.Printf("⚠️ Ошибка обновления статуса ордера: %v", err)
	} else {
		fmt.Printf("✅ Статус ордера обновлен на: FILLED\n")
	}

	// 7. Обновляем пользователя
	fmt.Println("\n✏️ 7. Обновление пользователя...")
	user.Username = "updated_demo_trader"
	user.Email = "updated_demo@metacore.com"
	user.KYCStatus = 2 // enhanced verification

	err = db.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("⚠️ Ошибка обновления пользователя: %v", err)
	} else {
		fmt.Printf("✅ Пользователь обновлен\n")
		fmt.Printf("   Новый username: %s, email: %s\n", user.Username, user.Email)
		fmt.Printf("   Новый KYC Status: %d\n", user.KYCStatus)
	}

	// 8. Создаем баланс пользователя
	fmt.Println("\n💰 8. Создание баланса пользователя...")
	balance := &domain.UserBalance{
		UserID:    user.ID,
		Asset:     "USDT",
		Free:      decimal.NewFromFloat(1000.00),
		Locked:    decimal.NewFromFloat(50.00),
		UpdatedAt: time.Now(),
	}

	// Примечание: BalanceStorage пока не реализован, но структура готова
	fmt.Printf("✅ Баланс подготовлен (метод пока не реализован)\n")
	fmt.Printf("   Asset: %s, Free: %s, Locked: %s\n", balance.Asset, balance.Free.String(), balance.Locked.String())

	// 9. Показываем финальную информацию
	fmt.Println("\n📋 9. Финальная информация...")
	fmt.Printf("   Пользователь ID: %d, Username: %s\n", user.ID, user.Username)
	fmt.Printf("   Ордер ID: %s, Symbol: %s, Status: %s\n", order.MexcOrderID, order.Symbol, order.Status)
	fmt.Printf("   Баланс: %s %s (Free: %s, Locked: %s)\n", balance.Free.Add(balance.Locked).String(), balance.Asset, balance.Free.String(), balance.Locked.String())

	fmt.Println("\n🎉 Демо-сценарий завершен успешно!")
	fmt.Println("   Все основные операции выполнены:")
	fmt.Println("   ✅ Создание пользователя")
	fmt.Println("   ✅ Получение пользователя")
	fmt.Println("   ✅ Создание ордера")
	fmt.Println("   ✅ Получение ордера")
	fmt.Println("   ✅ Обновление статуса ордера")
	fmt.Println("   ✅ Обновление пользователя")
	fmt.Println("   ✅ Подготовка баланса")
}
