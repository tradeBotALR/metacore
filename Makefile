# Makefile
.PHONY: lint lint-security lint-critic test build clean deps fix

# Основной линтинг
lint:
	golangci-lint run ./...

# Запуск всех линтеров включая дополнительные
lint-all:
	golangci-lint run --enable-all ./...

# Запуск линтеров с авто-фиксом
fix:
	golangci-lint run --fix ./...

# Быстрый линтинг
lint-fast:
	golangci-lint run --fast ./...

# Запуск gosec отдельно
gosec:
	gosec ./...

# Запуск gocritic отдельно
gocritic:
	gocritic check-project .

# Тесты
test:
	go test -v ./...

# Установка зависимостей
deps:
	go mod tidy

# Полная проверка
check: deps lint test

dep:
	go run ./postgres/cmd/ -db-url="postgres://postgres:postgres@localhost:5432/mexc_bot_db?sslmode=disable"
