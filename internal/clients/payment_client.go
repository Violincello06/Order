package clients

import (
	"context"
	"fmt"
	"time"
)

// PaymentRequestDTO data yang dikirim ke Payment Gateway saat inisiasi transaksi
type PaymentRequestDTO struct {
	OrderID       string  `json:"order_id"`
	GrossAmount   float64 `json:"gross_amount"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
}

// PaymentResponseDTO response dari Payment Gateway setelah transaksi dibuat
type PaymentResponseDTO struct {
	PaymentID     string    `json:"payment_id"`
	RedirectURL   string    `json:"redirect_url"`
	Status        string    `json:"status"` // e.g. "PENDING", "SUCCESS"
	ExpiryTime    time.Time `json:"expiry_time"`
}

// PaymentClient interface kontrak untuk berkomunikasi dengan Payment Gateway / API
type PaymentClient interface {
	CreateTransaction(ctx context.Context, req PaymentRequestDTO) (*PaymentResponseDTO, error)
	CheckTransactionStatus(ctx context.Context, paymentID string) (string, error)
}

// MockPaymentClient simulasi gateway pembayaran untuk tahap awal development
type MockPaymentClient struct{}

func NewMockPaymentClient() *MockPaymentClient {
	return &MockPaymentClient{}
}

func (c *MockPaymentClient) CreateTransaction(ctx context.Context, req PaymentRequestDTO) (*PaymentResponseDTO, error) {
	// Buat payment_id unik dan arahkan ke Sandbox Gateway simulator lokal yang aktif
	paymentID := fmt.Sprintf("PAY-%d", time.Now().UnixNano())
	redirectURL := fmt.Sprintf("/sandbox/pay/%s?order_id=%s&amount=%.0f&customer=%s", paymentID, req.OrderID, req.GrossAmount, req.CustomerName)

	return &PaymentResponseDTO{
		PaymentID:   paymentID,
		RedirectURL: redirectURL,
		Status:      "PENDING",
		ExpiryTime:  time.Now().Add(24 * time.Hour),
	}, nil
}

func (c *MockPaymentClient) CheckTransactionStatus(ctx context.Context, paymentID string) (string, error) {
	// Simulasi cek status transaksi ke payment gateway
	return "PENDING", nil
}

