package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"order/internal/models"
)

// Default Key Tokens (Bisa di-override lewat file .env / Environment Variables)
const (
	DefaultFrontendKey  = "ord_fe_live_9a7b6c5d4e3f2a1b9876543210fedcba"
	DefaultWebhookSecret = "whsec_pay_8841a0e9b62c451fbcda332156890abc"
	DefaultMasterKey     = "ord_master_admin_secret_99887766554433221100"
)

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// RequireFrontendAuth memvalidasi token API dari Frontend Client
func RequireFrontendAuth() echo.MiddlewareFunc {
	validKey := getEnvOrDefault("ORDER_API_KEY", DefaultFrontendKey)
	masterKey := getEnvOrDefault("MASTER_ADMIN_KEY", DefaultMasterKey)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Cek dari header X-API-Key
			apiKey := c.Request().Header.Get("X-API-Key")

			// 2. Jika tidak ada, cek dari Authorization: Bearer <token>
			if apiKey == "" {
				authHeader := c.Request().Header.Get("Authorization")
				if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
					apiKey = strings.TrimSpace(authHeader[7:])
				}
			}

			// Validasi token
			if apiKey == "" || (apiKey != validKey && apiKey != masterKey) {
				return c.JSON(http.StatusUnauthorized, models.ApiResponse{
					Success: false,
					Message: "Akses ditolak (Unauthorized): Sertakan API Key yang valid pada header 'X-API-Key' atau 'Authorization: Bearer <token>'",
				})
			}

			return next(c)
		}
	}
}

// RequireWebhookAuth memvalidasi token webhook dari Payment Gateway
func RequireWebhookAuth() echo.MiddlewareFunc {
	validSecret := getEnvOrDefault("PAYMENT_WEBHOOK_SECRET", DefaultWebhookSecret)
	masterKey := getEnvOrDefault("MASTER_ADMIN_KEY", DefaultMasterKey)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cek dari header X-Webhook-Token atau X-Payment-Signature
			webhookToken := c.Request().Header.Get("X-Webhook-Token")
			if webhookToken == "" {
				webhookToken = c.Request().Header.Get("X-Payment-Signature")
			}

			// Validasi token webhook
			if webhookToken == "" || (webhookToken != validSecret && webhookToken != masterKey) {
				return c.JSON(http.StatusUnauthorized, models.ApiResponse{
					Success: false,
					Message: "Akses ditolak (Unauthorized): Token webhook tidak valid. Sertakan token resmi pada header 'X-Webhook-Token'",
				})
			}

			return next(c)
		}
	}
}
