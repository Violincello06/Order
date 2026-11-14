package models

import (
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusCompleted OrderStatus = "COMPLETED"
)

// OrderItem merepresentasikan item produk di dalam pesanan
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

// Order merepresentasikan data lengkap sebuah pesanan
type Order struct {
	ID            string      `json:"id"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Items         []OrderItem `json:"items"`
	TotalAmount   float64     `json:"total_amount"`
	Status        OrderStatus `json:"status"`

	// Integrasi Payment API
	PaymentID      string `json:"payment_id,omitempty"`
	PaymentURL     string `json:"payment_url,omitempty"`
	PaymentStatus  string `json:"payment_status,omitempty"`
	PaymentMethod  string `json:"payment_method,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrderItemRequest input item saat membuat order
type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// CreateOrderRequest request payload untuk membuat order baru
type CreateOrderRequest struct {
	CustomerName  string                   `json:"customer_name"`
	CustomerEmail string                   `json:"customer_email"`
	Items         []CreateOrderItemRequest `json:"items"`
}

// PaymentWebhookRequest payload notifikasi webhook dari Payment Gateway
type PaymentWebhookRequest struct {
	OrderID       string  `json:"order_id"`
	PaymentID     string  `json:"payment_id"`
	TransactionStatus string `json:"transaction_status"` // e.g. "settlement", "success", "failed", "expire"
	PaymentMethod string  `json:"payment_method"`
	GrossAmount   float64 `json:"gross_amount"`
}

// ApiResponse standar response JSON untuk client
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

