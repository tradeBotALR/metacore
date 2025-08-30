package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samar/sup_bot/metacore/configs"

	_ "github.com/lib/pq"
)

// Deployer отвечает за развертывание и инициализацию БД
type Deployer struct {
	config configs.Config
}

// NewDeployer создает новый экземпляр Deployer
func NewDeployer(config configs.Config) *Deployer {
	return &Deployer{
		config: config,
	}
}

// DeploySchema развертывает схему БД
func (d *Deployer) DeploySchema(ctx context.Context) error {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.config.DB.Host, d.config.DB.Port, d.config.DB.User, d.config.DB.Password, d.config.DB.DBName, d.config.DB.SSLMode,
	)

	// Подключаемся к БД
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Connected to the database successfully.")

	// Читаем и выполняем схему
	if err := d.executeSchema(ctx, db); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Schema deployed successfully!")
	return nil
}

// executeSchema читает и выполняет SQL схему
func (d *Deployer) executeSchema(ctx context.Context, db *sql.DB) error {
	// Путь к файлу схемы
	schemaPath := "./postgres/schema.sql"

	// Читаем файл схемы
	sqlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	sql := string(sqlBytes)

	// Выполняем схему
	_, err = db.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// DeploySchemaWithTimeout развертывает схему с таймаутом
func (d *Deployer) DeploySchemaWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return d.DeploySchema(ctx)
}

// DeploySchemaDefault развертывает схему с таймаутом по умолчанию
func (d *Deployer) DeploySchemaDefault() error {
	return d.DeploySchemaWithTimeout(configs.DefaultConnectTimeout)
}
