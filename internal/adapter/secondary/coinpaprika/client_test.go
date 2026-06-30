package coinpaprika

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEthereumPrice_Success(t *testing.T) {
	mockResponse := `{"quotes": {"USD": {"price": 1850.55}}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	price, err := client.GetEthereumPrice(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPrice := 1850.55
	if price != expectedPrice {
		t.Errorf("expected price %f, got %f", expectedPrice, price)
	}
}

func TestGetEthereumPrice_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	_, err := client.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetEthereumPrice_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL, nil)
	_, err := client.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
