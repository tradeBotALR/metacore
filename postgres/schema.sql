-- Создание схемы базы данных

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     mexc_uid VARCHAR(255),
                                     username VARCHAR(255),
                                     email VARCHAR(255),
                                     mexc_api_key VARCHAR(255) NOT NULL,
                                     mexc_secret_key VARCHAR(255) NOT NULL,
                                     kyc_status SMALLINT,
                                     can_trade BOOLEAN,
                                     can_withdraw BOOLEAN,
                                     can_deposit BOOLEAN,
                                     account_type VARCHAR(50),
                                     permissions TEXT, -- JSON array
                                     last_account_sync TIMESTAMP,
                                     is_active BOOLEAN NOT NULL DEFAULT TRUE,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы ордеров
CREATE TABLE IF NOT EXISTS orders (
                                      id BIGSERIAL PRIMARY KEY,
                                      internal_id BIGINT UNIQUE, -- Для быстрых JOIN'ов
                                      user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      mexc_order_id VARCHAR(255) NOT NULL UNIQUE,
                                      symbol VARCHAR(20) NOT NULL,
                                      side VARCHAR(10) NOT NULL,
                                      type VARCHAR(20) NOT NULL,
                                      status VARCHAR(20) NOT NULL,
                                      price DECIMAL(30, 15),
                                      quantity DECIMAL(30, 15),
                                      quote_order_qty DECIMAL(30, 15),
                                      executed_quantity DECIMAL(30, 15),
                                      cummulative_quote_qty DECIMAL(30, 15),
                                      client_order_id VARCHAR(255),
                                      transact_time TIMESTAMP NOT NULL,
                                      created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы балансов пользователей
CREATE TABLE IF NOT EXISTS user_balances (
                                             id BIGSERIAL PRIMARY KEY,
                                             user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             asset VARCHAR(20) NOT NULL,
                                             free DECIMAL(30, 15) NOT NULL,
                                             locked DECIMAL(30, 15) NOT NULL,
                                             updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             UNIQUE(user_id, asset)
);

-- Создание таблицы сделок
CREATE TABLE IF NOT EXISTS trades (
                                      id BIGSERIAL PRIMARY KEY,
                                      user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                      mexc_trade_id VARCHAR(255) NOT NULL UNIQUE,
                                      order_id VARCHAR(255) NOT NULL, -- REFERENCES orders(mexc_order_id)
                                      symbol VARCHAR(20) NOT NULL,
                                      price DECIMAL(30, 15) NOT NULL,
                                      quantity DECIMAL(30, 15) NOT NULL,
                                      quote_quantity DECIMAL(30, 15) NOT NULL,
                                      commission DECIMAL(30, 15) NOT NULL,
                                      commission_asset VARCHAR(20) NOT NULL,
                                      trade_time TIMESTAMP NOT NULL,
                                      is_buyer BOOLEAN NOT NULL,
                                      is_maker BOOLEAN NOT NULL,
                                      created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы истории изменений статуса ордеров
CREATE TABLE IF NOT EXISTS order_updates (
                                             id BIGSERIAL PRIMARY KEY,
                                             user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             order_id VARCHAR(255) NOT NULL, -- REFERENCES orders(mexc_order_id)
                                             status VARCHAR(20) NOT NULL,
                                             executed_quantity DECIMAL(30, 15),
                                             cummulative_quote_qty DECIMAL(30, 15),
                                             update_time TIMESTAMP NOT NULL,
                                             raw_data JSONB
);