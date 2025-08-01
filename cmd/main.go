package main

import (
	"context"
	"github.com/shopspring/decimal"
	"log"
	"metacore/configs"
	"metacore/domain"
	"metacore/postgres"
	"metacore/storage/orders"
	"time"
)

func main() {
	cfg := configs.DefaultConfig()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Создаем экземпляр хранилища
	storage := orders.NewStorage(db.Pool())

	ctx := context.Background()

	// 3. Создаем ордер
	order := &domain.Order{
		// nolint : mnd
		UserID:      123,
		MexcOrderID: "some_mexc_order_id_12345",
		Symbol:      "BTCUSDT",
		Side:        "BUY",
		Type:        "LIMIT",
		Status:      "NEW",
		// nolint : mnd
		Price: decimal.NewFromFloat(50000.0),
		// nolint : mnd
		Quantity: decimal.NewFromFloat(0.001),
		// ... заполняем остальные поля
		TransactTime: time.Now(), // В миллисекундах
	}

	err = storage.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
	} else {
		log.Println("Order created successfully")
	}

	// 4. Получаем ордер
	retrievedOrder, err := storage.GetOrderByID(ctx, "some_mexc_order_id_12345")
	if err != nil {
		log.Printf("Failed to get order: %v", err)
	} else {
		log.Printf("Retrieved order: %+v", retrievedOrder)
	}

	// 5. Обновляем статус
	err = storage.UpdateOrderStatus(ctx, "some_mexc_order_id_12345", "FILLED")
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
	} else {
		log.Println("Order status updated successfully")
	}
	// 6. Удаляем ордер
	err = storage.DeleteOrderByID(ctx, "some_mexc_order_id_12345")

	if err != nil {
		log.Printf("Failed to delete order: %v", err)
	} else {
		log.Println("Order deleted successfully")
	}
}
