# Makefile
.PHONY: lint lint-security lint-critic test build clean deps fix

.PHONY: mock-gen
mock-gen:
	@echo "Generating mocks..."
	mockgen -destination=storage/mocks/full_storage_mock.go -package=mocks metacore/storage FullStorage
	@echo "Mocks generated successfully."

.PHONY: mock-gen-db
mock-gen-db:
	@echo "Database/sql mocks already created manually."
	@echo "If you need to regenerate them, use mockgen directly:"
	@echo "mockgen -destination=storage/mocks/db_interface_mock.go -package=mocks metacore/storage DBInterface"
	@echo "mockgen -destination=storage/mocks/sql_row_mock.go -package=mocks database/sql Row"
	@echo "mockgen -destination=storage/mocks/sql_result_mock.go -package=mocks database/sql Result"

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
	go test ./...

# Тесты только для UserStorage
test-users:
	go test ./postgres/internal/users/...

# Тесты только для OrderStorage  
test-orders:
	go test ./postgres/internal/orders/...

# Тесты только для BalancesStorage  
test-balances:
	go test ./postgres/internal/balances/...

# Тесты с покрытием
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Установка зависимостей
deps:
	go mod tidy

# Полная проверка
check: deps lint test

dep:
	go run ./postgres/cmd/ -db-url="postgres://postgres:postgres@localhost:5432/mexc_bot_db?sslmode=disable"
