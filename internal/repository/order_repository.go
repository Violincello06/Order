package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"order/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order tidak ditemukan")
)

// OrderRepository interface kontrak penyimpanan data order
type OrderRepository interface {
	Save(ctx context.Context, order *models.Order) error
	FindByID(ctx context.Context, id string) (*models.Order, error)
	FindAll(ctx context.Context) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
}

// MemoryOrderRepository implementasi in-memory repository untuk tahap awal
type MemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*models.Order),
	}
}

func (r *MemoryOrderRepository) Save(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return fmt.Errorf("order dengan ID %s sudah ada", order.ID)
	}

	r.orders[order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) FindByID(ctx context.Context, id string) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (r *MemoryOrderRepository) FindAll(ctx context.Context) ([]*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		result = append(result, order)
	}

	return result, nil
}

func (r *MemoryOrderRepository) Update(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return ErrOrderNotFound
	}

	r.orders[order.ID] = order
	return nil
}

