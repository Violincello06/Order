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

// Dynamic loader in-memory untuk me-mount App.vue langsung tanpa perlu file index.html di disk
const vueMountPage = `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Order Management | App.vue</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
  <script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
  <script src="https://unpkg.com/vue3-sfc-loader/dist/vue3-sfc-loader.js"></script>
</head>
<body class="bg-slate-50 m-0 p-0">
  <div id="app"></div>
  <script>
    const { loadModule } = window['vue3-sfc-loader'];
    const options = {
      moduleCache: { vue: Vue },
      async getFile(url) {
        const res = await fetch(url);
        if (!res.ok) throw new Error('Gagal membaca ' + url + ': ' + res.statusText);
        return { getContentData: asBinary => asBinary ? res.arrayBuffer() : res.text() };
      },
      addStyle(textContent) {
        const style = document.createElement('style');
        style.textContent = textContent;
        const ref = document.head.getElementsByTagName('style')[0] || null;
        document.head.insertBefore(style, ref);
      },
    };
    Vue.createApp(Vue.defineAsyncComponent(() => loadModule('/public/App.vue', options))).mount('#app');
  </script>
</body>
</html>`

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
	sandboxHandler := handler.NewSandboxHandler()

	// 4. Layani Tampilan App.vue (Tanpa file index.html di disk)
	renderAppVue := func(c echo.Context) error {
		return c.HTML(http.StatusOK, vueMountPage)
	}
	e.GET("/", renderAppVue)
	e.GET("/dashboard", renderAppVue)
	e.File("/App.vue", "public/App.vue")
	e.File("/docs", "public/docs.html")
	e.File("/openapi.json", "openapi.json")
	e.Static("/public", "public")

	// Layani Sandbox Payment Gateway Simulator Lokal
	sandboxHandler.RegisterRoutes(e)

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
