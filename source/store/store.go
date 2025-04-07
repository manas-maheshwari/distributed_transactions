package store

import (
	"context"
	"errors"
	"sync"

	"distributed_transactions/source/coordinator"
)

// Store represents a food store
type Store struct {
	*coordinator.BaseParticipant
	items     map[string]int // itemID -> quantity
	orders    map[string]map[string]int // transactionID -> (itemID -> quantity)
	mu        sync.RWMutex
}

// NewStore creates a new store instance
func NewStore() *Store {
	return &Store{
		BaseParticipant: coordinator.NewBaseParticipant(),
		items:          make(map[string]int),
		orders:         make(map[string]map[string]int),
	}
}

// AddItem adds an item to the store
func (s *Store) AddItem(itemID string, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[itemID] += quantity
}

// Prepare implements the Prepare method for the store
func (s *Store) Prepare(ctx context.Context, transactionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we have enough items for this transaction
	order, exists := s.orders[transactionID]
	if !exists {
		return errors.New("order not found")
	}

	for itemID, quantity := range order {
		if s.items[itemID] < quantity {
			return errors.New("insufficient items")
		}
	}

	// Reserve the items
	for itemID, quantity := range order {
		s.items[itemID] -= quantity
	}

	return s.BaseParticipant.Prepare(ctx, transactionID)
}

// Commit implements the Commit method for the store
func (s *Store) Commit(ctx context.Context, transactionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the order from our records
	delete(s.orders, transactionID)
	return s.BaseParticipant.Commit(ctx, transactionID)
}

// Abort implements the Abort method for the store
func (s *Store) Abort(ctx context.Context, transactionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return items to inventory
	order, exists := s.orders[transactionID]
	if exists {
		for itemID, quantity := range order {
			s.items[itemID] += quantity
		}
		delete(s.orders, transactionID)
	}

	return s.BaseParticipant.Abort(ctx, transactionID)
}

// PlaceOrder places an order for items
func (s *Store) PlaceOrder(ctx context.Context, transactionID string, items map[string]int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we have enough items
	for itemID, quantity := range items {
		if s.items[itemID] < quantity {
			return errors.New("insufficient items")
		}
	}

	// Record the order
	s.orders[transactionID] = items
	return nil
} 