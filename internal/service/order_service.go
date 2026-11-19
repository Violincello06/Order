package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"order/internal/clients"
	"order/internal/models"
	"order/internal/repository"
)

var (
	ErrInvalidInput       = errors.New("input data order tidak valid")
	ErrOrderCannotCancel  = errors.New("order tidak dapat dibatalkan karena status saat ini")
)

type OrderService interface {
	CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.Order, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	CancelOrder(ctx context.Context, id string) (*models.Order, error)
	ProcessPaymentWebhook(ctx context.Context, webhook models.PaymentWebhookRequest) error
}

type orderService struct {
	repo          repository.OrderRepository
	productClient clients.ProductClient
	paymentClient clients.PaymentClient
}

func NewOrderService(
	repo repository.OrderRepository,
	productClient clients.ProductClient,
	paymentClient clients.PaymentClient,
) OrderService {
	return &orderService{
		repo:          repo,
		productClient: productClient,
		paymentClient: paymentClient,
	}
}

// GenerateOrderID menghasilkan ID unik berformat ORD-YYYYMMDD-XXXX
func generateOrderID() string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	dateStr := time.Now().Format("20060102")
	return fmt.Sprintf("ORD-%s-%s", dateStr, strings.ToUpper(hex.EncodeToString(bytes)))
}

func (s *orderService) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.Order, error) {
	if req.CustomerName == "" || req.CustomerEmail == "" {
		return nil, fmt.Errorf("%w: nama dan email pemesan wajib diisi", ErrInvalidInput)
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("%w: minimal satu item produk harus dipesan", ErrInvalidInput)
	}

	// 1. Verifikasi data produk & stok via Product API Client
	var orderItems []models.OrderItem
	var totalAmount float64

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, fmt.Errorf("%w: quantity produk %s harus lebih dari 0", ErrInvalidInput, itemReq.ProductID)
		}

		// Ambil data produk terbaru dari Product Service
		product, err := s.productClient.GetProduct(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("validasi produk gagal: %w", err)
		}

		// Reservasi stok di Product Service
		if err := s.productClient.ReserveStock(ctx, itemReq.ProductID, itemReq.Quantity); err != nil {
			return nil, fmt.Errorf("gagal reservasi stok untuk produk '%s': %w", product.Name, err)
		}

		subtotal := product.Price * float64(itemReq.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    itemReq.Quantity,
			Subtotal:    subtotal,
			ImageURL:    product.ImageURL,
		})
	}

	orderID := generateOrderID()
	now := time.Now()

	// 2. Hubungi Payment API Client untuk inisiasi transaksi pembayaran
	paymentResp, err := s.paymentClient.CreateTransaction(ctx, clients.PaymentRequestDTO{
		OrderID:       orderID,
		GrossAmount:   totalAmount,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal inisiasi pembayaran di Payment Gateway: %w", err)
	}

	// 3. Simpan data Order ke Repository
	newOrder := &models.Order{
		ID:            orderID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Items:         orderItems,
		TotalAmount:   totalAmount,
		Status:        models.StatusPending,
		PaymentID:     paymentResp.PaymentID,
		PaymentURL:    paymentResp.RedirectURL,
		PaymentStatus: paymentResp.Status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Save(ctx, newOrder); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data order: %w", err)
	}

	return newOrder, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	return s.repo.FindAll(ctx)
}

func (s *orderService) CancelOrder(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Hanya order dengan status PENDING yang dapat dibatalkan
	if order.Status != models.StatusPending {
		return nil, fmt.Errorf("%w: status saat ini '%s'", ErrOrderCannotCancel, order.Status)
	}

	order.Status = models.StatusCancelled
	order.PaymentStatus = "CANCELLED"
	order.UpdatedAt = time.Now()

	// Rollback status barang di Product Service (reserved -> available)
	for _, item := range order.Items {
		_ = s.productClient.RollbackStock(ctx, item.ProductID)
	}

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) ProcessPaymentWebhook(ctx context.Context, webhook models.PaymentWebhookRequest) error {
	order, err := s.repo.FindByID(ctx, webhook.OrderID)
	if err != nil {
		return fmt.Errorf("order %s tidak ditemukan: %w", webhook.OrderID, err)
	}

	order.PaymentMethod = webhook.PaymentMethod
	order.PaymentStatus = strings.ToUpper(webhook.TransactionStatus)
	order.UpdatedAt = time.Now()

	// Update status order & sinkronisasi status barang ke Product Service
	switch strings.ToLower(webhook.TransactionStatus) {
	case "settlement", "success", "capture":
		order.Status = models.StatusPaid
		// Ubah status produk dari 'reserved' menjadi 'sold' di Product Service
		for _, item := range order.Items {
			_ = s.productClient.FinalizeSold(ctx, item.ProductID)
		}
	case "deny", "cancel", "expire", "failed":
		order.Status = models.StatusCancelled
		// Rollback status produk dari 'reserved' kembali ke 'available'
		for _, item := range order.Items {
			_ = s.productClient.RollbackStock(ctx, item.ProductID)
		}
	}

	return s.repo.Update(ctx, order)
}

