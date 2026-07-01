package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/steel-feel/prac/internal/port"
)

type priceService struct {
	repos []port.PriceRepository
}

// NewPriceService creates a new market price service.
func NewPriceService(repos ...port.PriceRepository) port.PriceService {
	return &priceService{repos: repos}
}

func (s *priceService) GetEthereumPrice(ctx context.Context) (float64, error) {
	if len(s.repos) == 0 {
		return 0, fmt.Errorf("no price providers configured")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var prices []float64

	for _, repo := range s.repos {
		wg.Add(1)
		go func(r port.PriceRepository) {
			defer wg.Done()
			price, err := r.GetEthereumPrice(ctx)
			if err == nil {
				mu.Lock()
				prices = append(prices, price)
				mu.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	if len(prices) == 0 {
		return 0, fmt.Errorf("all price providers failed")
	}

	var sum float64
	for _, p := range prices {
		sum += p
	}

	return sum / float64(len(prices)), nil
}
