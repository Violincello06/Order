package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"order/internal/models"
	"order/internal/repository"
	"order/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// RegisterRoutes mendaftarkan seluruh endpoint order ke Echo
func (h *OrderHandler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	// Endpoint Manajemen Order
	orders := api.Group("/orders")
	orders.POST("", h.CreateOrder)
	orders.GET("", h.GetAllOrders)
	orders.GET("/:id", h.GetOrderByID)
	orders.PUT("/:id/cancel", h.CancelOrder)
	orders.POST("/:id/pay", h.SimulatePayment) // Endpoint pembantu simulasi pembayaran

	// Endpoint Webhook Payment Gateway
	api.POST("/payments/webhook", h.HandlePaymentWebhook)
}

// CreateOrder handler untuk membuat pesanan baru
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req models.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Payload JSON tidak valid: " + err.Error(),
		})
	}

	order, err := h.orderService.CreateOrder(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "Order berhasil dibuat",
		Data:    order,
	})
}

// GetAllOrders handler untuk mengambil semua daftar pesanan
func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	orders, err := h.orderService.GetAllOrders(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: "Gagal mengambil data order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Daftar pesanan berhasil dimuat",
		Data:    orders,
	})
}

// GetOrderByID handler untuk mengambil detail spesifik suatu order
func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return c.JSON(http.StatusNotFound, models.ApiResponse{
				Success: false,
				Message: "Order tidak ditemukan",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Detail order berhasil diambil",
		Data:    order,
	})
}

// CancelOrder handler untuk membatalkan pesanan yang berstatus PENDING
func (h *OrderHandler) CancelOrder(c echo.Context) error {
	id := c.Param("id")
	order, err := h.orderService.CancelOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Order berhasil dibatalkan",
		Data:    order,
	})
}

// HandlePaymentWebhook handler untuk menerima callback / notifikasi dari Payment Gateway
func (h *OrderHandler) HandlePaymentWebhook(c echo.Context) error {
	var webhookReq models.PaymentWebhookRequest
	if err := c.Bind(&webhookReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "Payload webhook tidak valid: " + err.Error(),
		})
	}

	if err := h.orderService.ProcessPaymentWebhook(c.Request().Context(), webhookReq); err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Webhook pembayaran berhasil diproses",
	})
}

// SimulatePayment simulasi penyelesaian pembayaran untuk order tertentu
func (h *OrderHandler) SimulatePayment(c echo.Context) error {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ApiResponse{
			Success: false,
			Message: "Order tidak ditemukan",
		})
	}

	// Trigger webhook mock
	webhookReq := models.PaymentWebhookRequest{
		OrderID:           order.ID,
		PaymentID:         order.PaymentID,
		TransactionStatus: "settlement",
		PaymentMethod:     "qris",
		GrossAmount:       order.TotalAmount,
	}

	if err := h.orderService.ProcessPaymentWebhook(c.Request().Context(), webhookReq); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	updatedOrder, _ := h.orderService.GetOrderByID(c.Request().Context(), id)

	return c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "Pembayaran berhasil disimulasikan (status berubah menjadi PAID)",
		Data:    updatedOrder,
	})
}

