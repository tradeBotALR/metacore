package main

import (
	"context"
	"fmt"
	"log"

	"github.com/samar/sup_bot/metacore/configs"
	"github.com/samar/sup_bot/metacore/postgres"
)

func main() {
	fmt.Println("🔌 Тестирование подключения к базе данных...")

	// Загружаем конфигурацию
	cfg := configs.DefaultConfig()
	fmt.Printf("📋 Конфигурация:\n")
	fmt.Printf("   Host: %s\n", cfg.DB.Host)
	fmt.Printf("   Port: %d\n", cfg.DB.Port)
	fmt.Printf("   User: %s\n", cfg.DB.User)
	fmt.Printf("   Database: %s\n", cfg.DB.DBName)
	fmt.Printf("   SSL Mode: %s\n", cfg.DB.SSLMode)

	// Пытаемся подключиться к БД
	fmt.Println("\n🔗 Подключение к базе данных...")
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Printf("❌ Ошибка подключения к БД: %v", err)
		fmt.Println("\n💡 Возможные решения:")
		fmt.Println("   1. Убедитесь, что PostgreSQL запущен")
		fmt.Println("   2. Проверьте настройки подключения в configs/conf.go")
		fmt.Println("   3. Убедитесь, что база данных 'mexc_bot_db' существует")
		fmt.Println("   4. Проверьте права доступа пользователя 'postgres'")
		return
	}
	defer db.Close()

	fmt.Println("✅ Подключение к БД установлено успешно!")

	// Тестируем ping
	fmt.Println("\n🏓 Тестирование ping...")
	ctx := context.Background()
	err = db.Ping(ctx)
	if err != nil {
		log.Printf("❌ Ошибка ping: %v", err)
	} else {
		fmt.Println("✅ Ping успешен!")
	}

	fmt.Println("\n🎉 Тест подключения завершен успешно!")
	fmt.Println("   Теперь можно запускать полный демо-сценарий:")
	fmt.Println("   go run cmd/main.go")
}
