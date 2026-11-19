package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"order/internal/clients"
	"order/internal/handler"
	"order/internal/models"
	"order/internal/repository"
	"order/internal/service"
)

func main() {
	// 1. Inisialisasi Echo instance
	e := echo.New()

	// 2. Registrasi Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${uri} (${latency_human})\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// 3. Inisialisasi Komponen & Integrasi (Dependency Injection)
	// - Repositori Penyimpanan Order (In-memory, siap ditukar ke MySQL/GORM)
	orderRepo := repository.NewMemoryOrderRepository()

	// - Client Integrasi API Produk (mock / siap diganti real HTTP client)
	productClient := clients.NewMockProductClient()

	// - Client Integrasi API Payment (mock / siap diganti Midtrans, Xendit, Stripe, dll)
	paymentClient := clients.NewMockPaymentClient()

	// - Business Logic Service
	orderService := service.NewOrderService(orderRepo, productClient, paymentClient)

	// - HTTP Handlers
	orderHandler := handler.NewOrderHandler(orderService)

	// 4. Layani Web UI Frontend & API Documentation
	e.File("/", "public/index.html")
	e.File("/dashboard", "public/index.html")
	e.File("/docs", "public/docs.html")
	e.File("/openapi.json", "openapi.json")
	e.Static("/public", "public")

	// Health check & Info API
	e.GET("/api/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Message: "Order Management Service API is running",
			Data: map[string]interface{}{
				"version": "1.0.0",
				"docs_url": "/docs",
				"openapi_spec": "/openapi.json",
				"endpoints": []string{
					"POST   /api/v1/orders           (Buat order baru)",
					"GET    /api/v1/orders           (Ambil daftar order)",
					"GET    /api/v1/orders/:id       (Ambil detail order)",
					"PUT    /api/v1/orders/:id/cancel(Batalkan order)",
					"POST   /api/v1/orders/:id/pay   (Simulasi pembayaran sukses)",
					"POST   /api/v1/payments/webhook (Webhook notifikasi payment)",
				},
			},
		})
	})

	orderHandler.RegisterRoutes(e)

	// 5. Jalankan server pada port 1323
	e.Logger.Fatal(e.Start(":1323"))
}
