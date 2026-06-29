package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

// AccessLogMiddleware creates an access log for successful document requests.
func AccessLogMiddleware(repo port.AccessLogRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Call the next handler first
			err := next(c)

			id := c.Param("id")
			resp, errWrap := echo.UnwrapResponse(c.Response())
			status := 0
			if errWrap == nil {
				status = resp.Status
			}

			if err == nil && status == 200 && id != "" {
				ctx := c.Request().Context()
				
				txHash, _ := ctx.Value("x402_tx_hash").(string)
				payer, _ := ctx.Value("x402_payer").(string)

				log := &domain.AccessLog{
					DocumentID:        id,
					BlockchainAddress: payer,
					IPAddress:         c.RealIP(),
					UserAgent:         c.Request().UserAgent(),
					PaymentTxHash:     txHash,
				}

				// Save log asynchronously or synchronously depending on requirements
				// Doing it synchronously for simplicity here
				_ = repo.Save(ctx, log)
			}

			return err
		}
	}
}
