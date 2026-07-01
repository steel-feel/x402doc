package coinlore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetEthereumPrice_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ticker/" {
			t.Errorf("expected path /api/ticker/, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "80" {
			t.Errorf("expected id=80, got %s", r.URL.Query().Get("id"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"price_usd":"2950.45"}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	price, err := client.GetEthereumPrice(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 2950.45 {
		t.Errorf("expected price 2950.45, got %v", price)
	}
}

func TestClient_GetEthereumPrice_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_GetEthereumPrice_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"price_usd":"invalid"}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.GetEthereumPrice(context.Background())
	if err == nil {
		t.Fatal("expected error parsing float, got nil")
	}
}
