package delivery

import (
	"context"
	"errors"
	"sync"

	"distributed_transactions/source/coordinator"
)

// DeliveryAgent represents a delivery agent
type DeliveryAgent struct {
	ID        string
	Available bool
	Location  string
}

// Delivery represents the delivery service
type Delivery struct {
	*coordinator.BaseParticipant
	agents     map[string]*DeliveryAgent // agentID -> agent
	orders     map[string]string         // transactionID -> agentID
	mu         sync.RWMutex
}

// NewDelivery creates a new delivery service instance
func NewDelivery() *Delivery {
	return &Delivery{
		BaseParticipant: coordinator.NewBaseParticipant(),
		agents:         make(map[string]*DeliveryAgent),
		orders:         make(map[string]string),
	}
}

// AddAgent adds a delivery agent to the service
func (d *Delivery) AddAgent(agentID string, location string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.agents[agentID] = &DeliveryAgent{
		ID:        agentID,
		Available: true,
		Location:  location,
	}
}

// Prepare implements the Prepare method for the delivery service
func (d *Delivery) Prepare(ctx context.Context, transactionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	agentID, exists := d.orders[transactionID]
	if !exists {
		return errors.New("order not found")
	}

	agent, exists := d.agents[agentID]
	if !exists || !agent.Available {
		return errors.New("delivery agent not available")
	}

	// Mark agent as unavailable
	agent.Available = false
	return d.BaseParticipant.Prepare(ctx, transactionID)
}

// Commit implements the Commit method for the delivery service
func (d *Delivery) Commit(ctx context.Context, transactionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Remove the order from our records
	delete(d.orders, transactionID)
	return d.BaseParticipant.Commit(ctx, transactionID)
}

// Abort implements the Abort method for the delivery service
func (d *Delivery) Abort(ctx context.Context, transactionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Make agent available again
	agentID, exists := d.orders[transactionID]
	if exists {
		if agent, exists := d.agents[agentID]; exists {
			agent.Available = true
		}
		delete(d.orders, transactionID)
	}

	return d.BaseParticipant.Abort(ctx, transactionID)
}

// AssignAgent assigns a delivery agent to an order
func (d *Delivery) AssignAgent(ctx context.Context, transactionID string, location string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Find an available agent near the location
	for _, agent := range d.agents {
		if agent.Available {
			// TODO: Implement actual distance calculation
			// For now, we'll just assign the first available agent
			agent.Available = false
			d.orders[transactionID] = agent.ID
			return nil
		}
	}

	return errors.New("no delivery agents available")
} 