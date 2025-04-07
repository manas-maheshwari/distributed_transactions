package orders

import (
	"context"
	"errors"
	"sync"
	"time"

	"distributed_transactions/source/coordinator"
	"distributed_transactions/source/delivery"
	"distributed_transactions/source/store"
)

// Order represents a food order
type Order struct {
	ID        string
	Items     map[string]int // itemID -> quantity
	Location  string
	Status    string
	CreatedAt time.Time
}

// OrderService coordinates the 2PC transaction
type OrderService struct {
	coordinator *coordinator.Coordinator
	store       *store.Store
	delivery    *delivery.Delivery
	orders      map[string]*Order
	mu          sync.RWMutex
}

// NewOrderService creates a new order service instance
func NewOrderService(coordinator *coordinator.Coordinator, store *store.Store, delivery *delivery.Delivery) *OrderService {
	return &OrderService{
		coordinator: coordinator,
		store:       store,
		delivery:    delivery,
		orders:      make(map[string]*Order),
	}
}

// PlaceOrder places a new order and coordinates the 2PC transaction
func (o *OrderService) PlaceOrder(ctx context.Context, items map[string]int, location string) (*Order, error) {
	// Generate a unique transaction ID
	transactionID := generateTransactionID()

	// Create the order
	order := &Order{
		ID:        transactionID,
		Items:     items,
		Location:  location,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	o.mu.Lock()
	o.orders[transactionID] = order
	o.mu.Unlock()

	// Begin the transaction
	if _, err := o.coordinator.BeginTransaction(ctx, transactionID, 10*time.Minute); err != nil {
		return nil, err
	}

	// Phase 1: Prepare
	// 1. Place order in store
	if err := o.store.PlaceOrder(ctx, transactionID, items); err != nil {
		o.coordinator.Abort(ctx, transactionID)
		return nil, err
	}

	// 2. Assign delivery agent
	if err := o.delivery.AssignAgent(ctx, transactionID, location); err != nil {
		o.coordinator.Abort(ctx, transactionID)
		return nil, err
	}

	// 3. Prepare both participants
	participants := []string{"store", "delivery"}
	if err := o.coordinator.Prepare(ctx, transactionID, participants); err != nil {
		o.coordinator.Abort(ctx, transactionID)
		return nil, err
	}

	// Phase 2: Commit
	if err := o.coordinator.Commit(ctx, transactionID); err != nil {
		o.coordinator.Abort(ctx, transactionID)
		return nil, err
	}

	// Update order status
	order.Status = "confirmed"
	return order, nil
}

// GetOrder retrieves an order by ID
func (o *OrderService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, exists := o.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// generateTransactionID generates a unique transaction ID
func generateTransactionID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
} 