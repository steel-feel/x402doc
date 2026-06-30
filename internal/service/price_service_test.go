package service

import (
	"context"
	"errors"
	"testing"
)

type mockPriceRepository struct {
	price float64
	err   error
}

func (m *mockPriceRepository) GetEthereumPrice(ctx context.Context) (float64, error) {
	return m.price, m.err
}

func TestPriceService_GetEthereumPrice_Success(t *testing.T) {
	repo := &mockPriceRepository{
		price: 1920.45,
		err:   nil,
	}
	svc := NewPriceService(repo)

	price, err := svc.GetEthereumPrice(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPrice := 1920.45
	if price != expectedPrice {
		t.Errorf("expected price %f, got %f", expectedPrice, price)
	}
}

func TestPriceService_GetEthereumPrice_Error(t *testing.T) {
	repo := &mockPriceRepository{
		price: 0,
		err:   errors.New("db/network error"),
	}
	svc := NewPriceService(repo)

	_, err := svc.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
