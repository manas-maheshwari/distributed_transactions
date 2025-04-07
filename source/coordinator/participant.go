package coordinator

import (
	"context"
	"errors"
	"sync"
)

// Participant represents a service that can participate in a distributed transaction
type Participant interface {
	// Prepare prepares the participant for commit
	Prepare(ctx context.Context, transactionID string) error
	
	// Commit commits the transaction
	Commit(ctx context.Context, transactionID string) error
	
	// Abort aborts the transaction
	Abort(ctx context.Context, transactionID string) error
	
	// GetState returns the current state of the participant for the given transaction
	GetState(ctx context.Context, transactionID string) (TransactionState, error)
}

// BaseParticipant provides common functionality for participants
type BaseParticipant struct {
	states map[string]TransactionState
	mu     sync.RWMutex
}

// NewBaseParticipant creates a new base participant
func NewBaseParticipant() *BaseParticipant {
	return &BaseParticipant{
		states: make(map[string]TransactionState),
	}
}

// Prepare implements the Prepare method of Participant interface
func (p *BaseParticipant) Prepare(ctx context.Context, transactionID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if state, exists := p.states[transactionID]; exists && state != Initialized {
		return errors.New("invalid transaction state")
	}

	p.states[transactionID] = Prepared
	return nil
}

// Commit implements the Commit method of Participant interface
func (p *BaseParticipant) Commit(ctx context.Context, transactionID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if state, exists := p.states[transactionID]; !exists || state != Prepared {
		return errors.New("transaction not in prepared state")
	}

	p.states[transactionID] = Committed
	return nil
}

// Abort implements the Abort method of Participant interface
func (p *BaseParticipant) Abort(ctx context.Context, transactionID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.states[transactionID] = Aborted
	return nil
}

// GetState implements the GetState method of Participant interface
func (p *BaseParticipant) GetState(ctx context.Context, transactionID string) (TransactionState, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if state, exists := p.states[transactionID]; exists {
		return state, nil
	}

	return Initialized, errors.New("transaction not found")
} 