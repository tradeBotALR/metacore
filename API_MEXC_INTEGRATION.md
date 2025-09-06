# Интеграция с API MEXC (Spot v3)

Документ описывает, какие хранилищные методы реализованы в `metacore`, как они соответствуют объектам и эндпоинтам MEXC, и что ещё требуется добавить для полной поддержки.

## Реализованные хранилища и методы

- Пользователи (users)
  - CreateUser, GetUserByID, GetUserByMexcUID, GetUserByTelegramID, GetAllUsers, UpdateUser, DeleteUser
- Ордеры (orders)
  - CreateOrder, GetOrderByID, UpdateOrderStatus, DeleteOrderByID
- Сделки (trades)
  - CreateTrade, GetTradeByID, GetUserTrades(filters: symbol, startTime, endTime, limit, offset)
- Балансы (user_balances)
  - UpdateBalance(applied bool), GetBalance, GetUserBalances, UpdateUserBalances(транзакционно)

Все числовые денежные поля — `shopspring/decimal`. Времена — `time.Time` (сервер должен конвертировать миллисекунды UNIX из MEXC).

## Маппинг на объекты и эндпоинты MEXC

- Балансы: GET `/api/v3/account`
  - Server получает `balances[]` и вызывает `UpdateUserBalances` для upsert всех активов пользователя.
- Ордеры: POST/GET/DELETE `/api/v3/order`
  - POST (создание) → после ответа биржи сохраняем в БД через `CreateOrder`.
  - GET (запрос состояния) → можно прочитать из БД `GetOrderByID` (источник правды — биржа, БД кэширующая).
  - DELETE (отмена) → при успешной отмене вызвать `UpdateOrderStatus` или `DeleteOrderByID` (в зависимости от вашей политики хранения).
- Сделки: GET `/api/v3/myTrades`
  - Загруженные с биржи сделки сохраняем `CreateTrade`. Для выборки в UI — `GetUserTrades`.

## Поведение методов (кратко)

- UpdateUserBalances: транзакция, upsert по `(user_id, asset)`, единый `updated_at`, удаление нулевых балансов (free=0 и locked=0), защищает от частичных обновлений.
- UpdateBalance: upsert одной записи; возвращает `applied=true` только если значения изменились.
- UpdateOrderStatus: обновляет статус ордера; возвращает ошибку, если строк не затронуто.
- CreateTrade: сохраняет сделку; уникальность по `mexc_trade_id`.
- GetUserTrades: сортировка по `trade_time desc, id desc` + фильтры и пагинация.

## Что ещё требуется для полной поддержки Spot v3

- Расширение модели ордера под расширенные типы MEXC (при необходимости):
  - timeInForce, stopPrice, icebergQty и др. (эндпоинт `/api/v3/order`)
- Методы для витрин/выборок ордеров:
  - GetUserOrders(userID, filters: symbol, status, time range, limit/offset)
  - GetOpenOrders(userID, symbol?) — маппится на GET `/api/v3/openOrders`
- История изменений ордеров (таблица `order_updates` есть, методов нет):
  - AppendOrderUpdate(update), GetOrderUpdates(orderID, range)
- Методы для сделок по ордеру:
  - GetTradesByOrderID(orderID)
- Индексы БД (рекомендации):
  - `CREATE INDEX idx_user_balances_user_asset ON user_balances(user_id, asset);`
  - `CREATE INDEX idx_trades_user_time ON trades(user_id, trade_time DESC);`
  - `CREATE INDEX idx_trades_user_symbol ON trades(user_id, symbol);`
- Безопасность ключей MEXC:
  - Хранение зашифрованно (KMS/Vault) и ограниченный доступ в приложении.

## Практический сценарий

Пример сценария находится в `examples/order_scenario.go` и вызывается из `cmd/main.go`:

1) Создать пользователя → 2) Создать ордер → 3) Прочитать ордер → 4) Обновить статус (FILLED) → 5) Создать сделку → 6) Обновить балансы (USDT↓, BTC↑) → 7) Удалить ордер.

Запуск:
```bash
go run ./cmd
```

## Примечания по интеграции (Server слой)

- Сервер подписывает запросы к MEXC (timestamp, signature) и конвертирует время (ms → time.Time).
- После ответов биржи сервер вызывает методы `metacore` для записи/чтения состояния.
- `metacore` не делает сетевые вызовы — только работа с БД, что упрощает тестирование и изоляцию.
