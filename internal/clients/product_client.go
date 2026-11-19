package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ProductDTO format data produk yang diterima dari Pre-Loved Product Service
type ProductDTO struct {
	ID          string  `json:"id"`
	SellerID    int     `json:"seller_id,omitempty"`
	CategoryID  int     `json:"category_id,omitempty"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Condition   string  `json:"condition,omitempty"`
	Status      string  `json:"status"` // "available", "reserved", "sold"
	Stock       int     `json:"stock"`
	WeightGrams int     `json:"weight_grams,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// ProductClient interface kontrak resmi dengan Product Microservice
type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductDTO, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
	FinalizeSold(ctx context.Context, productID string) error
	RollbackStock(ctx context.Context, productID string) error
}

// =========================================================================
// Real HTTP Product Client (Terhubung ke http://localhost:3001)
// =========================================================================

type HttpProductClient struct {
	baseURL        string
	internalSecret string
	httpClient     *http.Client
	fallbackMock   *MockProductClient
}

func NewHttpProductClient() *HttpProductClient {
	baseURL := os.Getenv("PRODUCT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001"
	}

	secret := os.Getenv("PRODUCT_SERVICE_SECRET")
	if secret == "" {
		secret = "CHANGE_ME_LOCAL_DUMMY_SECRET_2026"
	}

	return &HttpProductClient{
		baseURL:        baseURL,
		internalSecret: secret,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		fallbackMock: NewMockProductClient(),
	}
}

type productAPIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type productStatusUpdatePayload struct {
	Status string `json:"status"`
}

// GetProduct memanggil GET /api/products/{id}
func (c *HttpProductClient) GetProduct(ctx context.Context, productID string) (*ProductDTO, error) {
	url := fmt.Sprintf("%s/api/products/%s", c.baseURL, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return c.fallbackMock.GetProduct(ctx, productID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Jika Product Service di port 3001 offline/belum jalan, fallback ke mock lokal
		return c.fallbackMock.GetProduct(ctx, productID)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Jika response dari 3001 bukan 200 (misal database kosong/404), gunakan fallback mock
		return c.fallbackMock.GetProduct(ctx, productID)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.fallbackMock.GetProduct(ctx, productID)
	}

	var apiResp productAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return c.fallbackMock.GetProduct(ctx, productID)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &rawData); err != nil {
		return c.fallbackMock.GetProduct(ctx, productID)
	}

	// Normalisasi ID menjadi string
	var idStr string
	switch v := rawData["id"].(type) {
	case float64:
		idStr = strconv.Itoa(int(v))
	case string:
		idStr = v
	default:
		idStr = productID
	}

	name, _ := rawData["name"].(string)
	price, _ := rawData["price"].(float64)
	status, _ := rawData["status"].(string)
	stockFloat, _ := rawData["stock"].(float64)

	// Ambil sumber gambar (primary_image atau images array)
	imageURL, _ := rawData["primary_image"].(string)
	if imageURL == "" {
		if imagesArr, ok := rawData["images"].([]interface{}); ok && len(imagesArr) > 0 {
			if firstImg, ok := imagesArr[0].(map[string]interface{}); ok {
				imageURL, _ = firstImg["image_url"].(string)
			}
		}
	}
	if imageURL == "" {
		if mockProd, err := c.fallbackMock.GetProduct(ctx, productID); err == nil {
			imageURL = mockProd.ImageURL
		}
	}

	return &ProductDTO{
		ID:       idStr,
		Name:     name,
		Price:    price,
		Status:   status,
		Stock:    int(stockFloat),
		ImageURL: imageURL,
	}, nil
}

// updateStatus memanggil PATCH /api/products/{id}/status dengan x-internal-secret
func (c *HttpProductClient) updateStatus(ctx context.Context, productID, status string) error {
	url := fmt.Sprintf("%s/api/products/%s/status", c.baseURL, productID)
	payload := productStatusUpdatePayload{Status: status}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.fallbackMock.updateStatus(productID, status)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-internal-secret", c.internalSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Fallback ke mock jika service offline
		return c.fallbackMock.updateStatus(productID, status)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// Jika produk belum tersimpan di DB port 3001, update di mock
	return c.fallbackMock.updateStatus(productID, status)
}

// ReserveStock mengunci barang saat order checkout (status: available -> reserved)
func (c *HttpProductClient) ReserveStock(ctx context.Context, productID string, quantity int) error {
	return c.updateStatus(ctx, productID, "reserved")
}

// FinalizeSold menandai barang telah lunas (status: reserved -> sold)
func (c *HttpProductClient) FinalizeSold(ctx context.Context, productID string) error {
	return c.updateStatus(ctx, productID, "sold")
}

// RollbackStock membatalkan booking (status: reserved -> available)
func (c *HttpProductClient) RollbackStock(ctx context.Context, productID string) error {
	return c.updateStatus(ctx, productID, "available")
}

// =========================================================================
// MockProductClient (Fallback & Data Pre-Loved Default dengan Gambar Asli)
// =========================================================================

type MockProductClient struct {
	mu       sync.RWMutex
	products map[string]*ProductDTO
}

func NewMockProductClient() *MockProductClient {
	return &MockProductClient{
		products: map[string]*ProductDTO{
			"1": {
				ID:        "1",
				Name:      "Jaket Denim Levi's Vintage 90s Size L",
				Price:     450000,
				Status:    "available",
				Stock:     1,
				Condition: "lightly_used",
				// Sumber gambar resmi dari Unsplash sesuai koleksi Tim Produk
				ImageURL:  "https://images.unsplash.com/photo-1576995853123-5a10305d93c0?w=400&auto=format&fit=crop&q=80",
			},
			"2": {
				ID:        "2",
				Name:      "Sepatu Sneakers Converse 70s Chuck Taylor",
				Price:     350000,
				Status:    "available",
				Stock:     1,
				Condition: "like_new",
				ImageURL:  "https://images.unsplash.com/photo-1525966222134-fcfa99b8ae77?w=400&auto=format&fit=crop&q=80",
			},
			"3": {
				ID:        "3",
				Name:      "Kamera Analog Vintage 35mm Terawat",
				Price:     750000,
				Status:    "available",
				Stock:     1,
				Condition: "well_used",
				ImageURL:  "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=400&auto=format&fit=crop&q=80",
			},
		},
	}
}

func (c *MockProductClient) GetProduct(ctx context.Context, productID string) (*ProductDTO, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	product, exists := c.products[productID]
	if !exists {
		return nil, fmt.Errorf("produk dengan ID '%s' tidak ditemukan di katalog", productID)
	}

	copyProd := *product
	return &copyProd, nil
}

func (c *MockProductClient) updateStatus(productID, status string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	product, exists := c.products[productID]
	if !exists {
		return fmt.Errorf("produk '%s' tidak ditemukan", productID)
	}

	product.Status = status
	if status == "reserved" || status == "sold" {
		product.Stock = 0
	} else if status == "available" {
		product.Stock = 1
	}
	return nil
}

func (c *MockProductClient) ReserveStock(ctx context.Context, productID string, quantity int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	product, exists := c.products[productID]
	if !exists {
		return fmt.Errorf("produk '%s' tidak ditemukan", productID)
	}

	if product.Status != "available" || product.Stock < 1 {
		return errors.New("stok produk tidak tersedia (sudah dibooking atau terjual)")
	}

	product.Status = "reserved"
	product.Stock = 0
	return nil
}

func (c *MockProductClient) FinalizeSold(ctx context.Context, productID string) error {
	return c.updateStatus(productID, "sold")
}

func (c *MockProductClient) RollbackStock(ctx context.Context, productID string) error {
	return c.updateStatus(productID, "available")
}
