package trades

import (
	"context"
	"metacore/domain"
	"metacore/storage"
)

// TradeStorage реализует интерфейс TradeStorage.
type TradeStorage struct {
	// TODO: Реализовать методы для работы со сделками
}

// NewTradeStorage создает новый экземпляр TradeStorage.
func NewTradeStorage() *TradeStorage {
	return &TradeStorage{}
}

// CreateTrade создает новую сделку.
func (s *TradeStorage) CreateTrade(ctx context.Context, trade *domain.Trade) error {
	// TODO: Реализовать создание сделки
	return nil
}

// GetTradeByID получает сделку по MEXC Trade ID.
func (s *TradeStorage) GetTradeByID(ctx context.Context, mexcTradeID string) (*domain.Trade, error) {
	// TODO: Реализовать получение сделки
	return nil, nil
}

// Ensure TradeStorage implements TradeStorage interface
var _ storage.TradeStorage = (*TradeStorage)(nil)
