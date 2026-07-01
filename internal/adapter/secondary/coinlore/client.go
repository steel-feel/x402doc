package coinlore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/steel-feel/prac/internal/port"
)

type coinloreResponse []struct {
	PriceUSD string `json:"price_usd"`
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Coinlore API client.
func NewClient(baseURL string, httpClient *http.Client) port.PriceRepository {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *client) GetEthereumPrice(ctx context.Context) (float64, error) {
	url := fmt.Sprintf("%s/api/ticker/?id=80", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data coinloreResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no data returned from coinlore")
	}

	price, err := strconv.ParseFloat(data[0].PriceUSD, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price_usd: %w", err)
	}

	return price, nil
}
