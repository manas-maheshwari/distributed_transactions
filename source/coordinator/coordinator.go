package coordinator

import (
	"context"
	"errors"
	"sync"
	"time"
)

// TransactionState represents the state of a distributed transaction
type TransactionState int

const (
	Initialized TransactionState = iota
	Prepared
	Committed
	Aborted
)

// Transaction represents a distributed transaction
type Transaction struct {
	ID        string
	State     TransactionState
	CreatedAt time.Time
	Timeout   time.Duration
	mu        sync.RWMutex
}

// Coordinator manages distributed transactions
type Coordinator struct {
	transactions map[string]*Transaction
	mu           sync.RWMutex
}

// NewCoordinator creates a new coordinator instance
func NewCoordinator() *Coordinator {
	return &Coordinator{
		transactions: make(map[string]*Transaction),
	}
}

// BeginTransaction starts a new distributed transaction
func (c *Coordinator) BeginTransaction(ctx context.Context, transactionID string, timeout time.Duration) (*Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.transactions[transactionID]; exists {
		return nil, errors.New("transaction already exists")
	}

	tx := &Transaction{
		ID:        transactionID,
		State:     Initialized,
		CreatedAt: time.Now(),
		Timeout:   timeout,
	}

	c.transactions[transactionID] = tx
	return tx, nil
}

// Prepare asks participants to prepare for commit
func (c *Coordinator) Prepare(ctx context.Context, transactionID string, participants []string) error {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if !exists {
		return errors.New("transaction not found")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != Initialized {
		return errors.New("invalid transaction state")
	}

	// TODO: Implement actual prepare calls to participants
	// For now, we'll simulate successful preparation
	tx.State = Prepared
	return nil
}

// Commit finalizes the transaction
func (c *Coordinator) Commit(ctx context.Context, transactionID string) error {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if !exists {
		return errors.New("transaction not found")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.State != Prepared {
		return errors.New("transaction not in prepared state")
	}

	// TODO: Implement actual commit calls to participants
	tx.State = Committed
	return nil
}

// Abort cancels the transaction
func (c *Coordinator) Abort(ctx context.Context, transactionID string) error {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if !exists {
		return errors.New("transaction not found")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	// TODO: Implement actual abort calls to participants
	tx.State = Aborted
	return nil
}

// GetTransactionState returns the current state of a transaction
func (c *Coordinator) GetTransactionState(ctx context.Context, transactionID string) (TransactionState, error) {
	c.mu.RLock()
	tx, exists := c.transactions[transactionID]
	c.mu.RUnlock()

	if !exists {
		return Initialized, errors.New("transaction not found")
	}

	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return tx.State, nil
} 