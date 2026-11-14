package clients

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ProductDTO format data produk yang diterima dari API Produk
type ProductDTO struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ProductClient interface kontrak untuk berkomunikasi dengan Product Service / API
type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductDTO, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
}

// MockProductClient implementasi awal berbasis memory (simulasi API Produk)
type MockProductClient struct {
	mu       sync.RWMutex
	products map[string]*ProductDTO
}

// NewMockProductClient membuat mock client dengan data produk awal untuk pengujian
func NewMockProductClient() *MockProductClient {
	return &MockProductClient{
		products: map[string]*ProductDTO{
			"PROD-001": {ID: "PROD-001", Name: "Kopi Arabika 250g", Price: 75000, Stock: 50},
			"PROD-002": {ID: "PROD-002", Name: "Syrup Vanilla 500ml", Price: 45000, Stock: 20},
			"PROD-003": {ID: "PROD-003", Name: "Paper Cup 12oz (Pack of 50)", Price: 35000, Stock: 100},
		},
	}
}

func (c *MockProductClient) GetProduct(ctx context.Context, productID string) (*ProductDTO, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	product, exists := c.products[productID]
	if !exists {
		return nil, fmt.Errorf("produk dengan ID '%s' tidak ditemukan di Product Service", productID)
	}

	// Salinan data agar aman dari modifikasi eksternal
	copyProd := *product
	return &copyProd, nil
}

func (c *MockProductClient) ReserveStock(ctx context.Context, productID string, quantity int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	product, exists := c.products[productID]
	if !exists {
		return fmt.Errorf("produk '%s' tidak ditemukan", productID)
	}

	if product.Stock < quantity {
		return errors.New("stok produk tidak mencukupi")
	}

	product.Stock -= quantity
	return nil
}

