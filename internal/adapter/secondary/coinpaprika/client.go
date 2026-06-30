package coinpaprika

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/steel-feel/prac/internal/port"
)

type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new CoinPaprika HTTP API client.
func NewClient(baseURL string, httpClient *http.Client) port.PriceRepository {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

type tickerResponse struct {
	Quotes struct {
		USD struct {
			Price float64 `json:"price"`
		} `json:"USD"`
	} `json:"quotes"`
}

func (c *client) GetEthereumPrice(ctx context.Context) (float64, error) {
	url := fmt.Sprintf("%s/v1/tickers/eth-ethereum", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data tickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Quotes.USD.Price, nil
}
