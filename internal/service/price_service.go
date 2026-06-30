package service

import (
	"context"

	"github.com/steel-feel/prac/internal/port"
)

type priceService struct {
	repo port.PriceRepository
}

// NewPriceService creates a new market price service.
func NewPriceService(repo port.PriceRepository) port.PriceService {
	return &priceService{repo: repo}
}

func (s *priceService) GetEthereumPrice(ctx context.Context) (float64, error) {
	return s.repo.GetEthereumPrice(ctx)
}
