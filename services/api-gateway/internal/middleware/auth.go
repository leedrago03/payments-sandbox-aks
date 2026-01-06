package middleware

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type ValidationResponse struct {
	Valid      bool   `json:"valid"`
	MerchantID string `json:"merchant_id"`
	Error      string `json:"error"`
}

func APIKeyAuth(merchantServiceURL string) fiber.Handler {
	client := &fasthttp.Client{
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	}

	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-API-Key header",
			})
		}

		// TODO: Add Redis caching here

		// Call Merchant Service to verify
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.SetRequestURI(merchantServiceURL + "/internal/api-keys/verify")
		req.Header.SetMethod(fasthttp.MethodPost)
		req.Header.SetContentType("application/json")
		
		body := map[string]string{"api_key": apiKey}
		jsonBody, _ := json.Marshal(body)
		req.SetBody(jsonBody)

		if err := client.Do(req, resp); err != nil {
			// Fail open or closed? Closed.
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Auth service unavailable",
			})
		}

		if resp.StatusCode() != fiber.StatusOK {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		}

		var validation ValidationResponse
		if err := json.Unmarshal(resp.Body(), &validation); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Auth response error",
			})
		}

		if !validation.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API Key",
			})
		}

		// Inject Merchant ID into headers for downstream services
		c.Request().Header.Set("X-Merchant-ID", validation.MerchantID)

		return c.Next()
	}
}
