package configs

import "time"

const (
	DefaultConnectTimeout    = 10 * time.Second // Время ожидания подключения
	DefaultMaxConnIdleTime   = 30 * time.Minute // Максимальное время простоя соединения
	DefaultMaxConnLifetime   = 1 * time.Hour    // Максимальное время жизни соединения
	DefaultHealthCheckPeriod = 1 * time.Minute  // Период проверки состояния соединений
	DefaultMaxConns          = int32(20)        // Максимальное количество соединений в пуле
	DefaultMinConns          = int32(5)         // Минимальное количество соединений в пуле

	DefaultPort = 5432
	DefaultHost = "localhost"
	DefaultUser = "postgres"
	DevSchema   = "mexc_bot_db"
)

// DBConfig описывает конфигурацию подключения к базе данных.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// PoolConfig описывает конфигурацию пула соединений.
type PoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

// Config объединяет всю конфигурацию БД.
type Config struct {
	DB   DBConfig
	Pool PoolConfig
}

func DefaultConfig() Config {
	return Config{
		DB: DBConfig{
			Host:     DefaultHost,
			Port:     DefaultPort,
			User:     DefaultUser, // Значения по умолчанию из docker-compose.yml
			Password: DefaultUser,
			DBName:   DevSchema,
			SSLMode:  "disable",
		},
		Pool: PoolConfig{
			MaxConns:          DefaultMaxConns,
			MinConns:          DefaultMinConns,
			MaxConnLifetime:   DefaultMaxConnLifetime,
			MaxConnIdleTime:   DefaultMaxConnIdleTime,
			HealthCheckPeriod: DefaultHealthCheckPeriod,
			ConnectTimeout:    DefaultConnectTimeout,
		},
	}
}
