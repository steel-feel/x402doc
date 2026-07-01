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
	repo1 := &mockPriceRepository{price: 1500.00, err: nil}
	repo2 := &mockPriceRepository{price: 1600.00, err: nil}

	svc := NewPriceService(repo1, repo2)
	price, err := svc.GetEthereumPrice(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if price != 1550.00 {
		t.Errorf("expected price 1550.00, got %v", price)
	}
}

func TestPriceService_GetEthereumPrice_PartialFailure(t *testing.T) {
	repo1 := &mockPriceRepository{price: 1500.00, err: nil}
	repo2 := &mockPriceRepository{price: 0, err: errors.New("provider failed")}

	svc := NewPriceService(repo1, repo2)
	price, err := svc.GetEthereumPrice(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if price != 1500.00 {
		t.Errorf("expected price 1500.00, got %v", price)
	}
}

func TestPriceService_GetEthereumPrice_TotalFailure(t *testing.T) {
	repo1 := &mockPriceRepository{price: 0, err: errors.New("provider 1 failed")}
	repo2 := &mockPriceRepository{price: 0, err: errors.New("provider 2 failed")}

	svc := NewPriceService(repo1, repo2)
	_, err := svc.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPriceService_GetEthereumPrice_NoProviders(t *testing.T) {
	svc := NewPriceService()
	_, err := svc.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
