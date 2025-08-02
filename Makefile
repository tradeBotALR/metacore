# Makefile
.PHONY: lint lint-security lint-critic test build clean deps fix

.PHONY: mock-gen
mock-gen:
	@echo "Generating mocks..."
	mockgen -destination=storage/mocks/full_storage_mock.go -package=mocks metacore/storage FullStorage
	@echo "Mocks generated successfully."

.PHONY: mock-gen-pgx
mock-gen-pgx:
	@echo "Generating pgxpool mocks..."
	mockgen -destination=storage/mocks/pgxpool_iface_mock.go -package=mocks metacore/storage PgxPoolIface
	@echo "PGX Pool mocks generated successfully."

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
