package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/internal/adapter/secondary/facilitator"
	"github.com/steel-feel/prac/internal/port"
)

// X402Middleware returns an Echo middleware that enforces x402 payments.
func X402Middleware(docSvc port.DocumentService, fac *facilitator.Facilitator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			id := c.Param("id")
			doc, err := docSvc.GetDocument(c.Request().Context(), id)
			if err != nil || doc == nil {
				// If doc isn't found, let the actual handler return 404
				return next(c)
			}

			// Check for PAYMENT-SIGNATURE header
			sig := c.Request().Header.Get("PAYMENT-SIGNATURE")
			if sig == "" {
				// Send 402 with PAYMENT-REQUIRED header
				req := map[string]any{
					"x402Version": 1,
					"accepts": []map[string]any{
						{
							"scheme":            "exact",
							"network":           "tempo-testnet", // Assuming Tempo testnet
							"maxAmountRequired": doc.PriceUSD,
							// "payTo" and "asset" would come from facilitator config ideally
						},
					},
				}
				reqBytes, _ := json.Marshal(req)
				c.Response().Header().Set("PAYMENT-REQUIRED", base64.StdEncoding.EncodeToString(reqBytes))
				return c.JSON(http.StatusPaymentRequired, map[string]string{"error": "payment required"})
			}

			// Verify payment
			res, err := fac.VerifyPayment(c.Request().Context(), sig, doc.PriceUSD)
			if err != nil || !res.Success {
				return c.JSON(http.StatusPaymentRequired, map[string]string{"error": "payment verification failed"})
			}

			// Store payment info in context for access logging
			ctx := context.WithValue(c.Request().Context(), "x402_tx_hash", res.TxHash)
			ctx = context.WithValue(ctx, "x402_payer", res.Payer)
			c.SetRequest(c.Request().WithContext(ctx))

			// Add PAYMENT-RESPONSE header
			resp := map[string]any{
				"txHash": res.TxHash,
				"status": "settled",
			}
			respBytes, _ := json.Marshal(resp)
			c.Response().Header().Set("PAYMENT-RESPONSE", base64.StdEncoding.EncodeToString(respBytes))

			return next(c)
		}
	}
}
