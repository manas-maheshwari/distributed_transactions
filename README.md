# Distributed Transactions for Food Delivery System

A Two-Phase Commit (2PC) implementation for a food delivery system that guarantees 10-minute delivery by ensuring both food items and delivery agents are available before confirming an order.

## Overview

This project implements a distributed transaction system using the Two-Phase Commit protocol to coordinate between multiple services:
- Order Service: Coordinates the overall transaction
- Store Service: Manages food item inventory
- Delivery Service: Manages delivery agent availability
- Coordinator: Implements the 2PC protocol

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Order      │     │  Store      │     │  Delivery   │
│  Service    │◄───►│  Service    │◄───►│  Service    │
└─────────────┘     └─────────────┘     └─────────────┘
        ▲                  ▲                  ▲
        │                  │                  │
        ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────┐
│                    Coordinator                      │
│              (Two-Phase Commit)                     │
└─────────────────────────────────────────────────────┘
```

### Components

1. **Order Service**
   - Manages order lifecycle
   - Coordinates between store and delivery services
   - Generates unique transaction IDs
   - Maintains order status

2. **Store Service**
   - Manages food item inventory
   - Locks items during transaction
   - Releases items if transaction fails
   - Implements participant interface for 2PC

3. **Delivery Service**
   - Manages delivery agent availability
   - Assigns agents to orders
   - Releases agents if transaction fails
   - Implements participant interface for 2PC

4. **Coordinator**
   - Implements Two-Phase Commit protocol
   - Manages transaction states
   - Handles transaction timeouts
   - Coordinates between participants

## Two-Phase Commit Protocol

1. **Prepare Phase**
   - Order service initiates transaction
   - Store service checks and reserves items
   - Delivery service checks and reserves agent
   - All participants must agree to proceed

2. **Commit Phase**
   - If all participants agree, transaction is committed
   - Items are deducted from inventory
   - Agent is assigned to delivery
   - Order is confirmed

3. **Abort Phase**
   - If any participant fails to prepare, transaction is aborted
   - Reserved items are returned to inventory
   - Reserved agent is made available again
   - Order is marked as failed

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/distributed_transactions.git
cd distributed_transactions
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run source/main.go
```

## Usage Example

```go
// Initialize services
coord := coordinator.NewCoordinator()
storeService := store.NewStore()
deliveryService := delivery.NewDelivery()
orderService := orders.NewOrderService(coord, storeService, deliveryService)

// Add initial data
storeService.AddItem("pizza", 5)
storeService.AddItem("burger", 10)
deliveryService.AddAgent("agent1", "location1")

// Place an order
items := map[string]int{
    "pizza": 2,
    "burger": 1,
}
order, err := orderService.PlaceOrder(ctx, items, "customer_location")
```

## Features

- Distributed transaction management
- Resource locking and release
- Transaction timeout handling
- Concurrent order processing
- Consistent state management
- 10-minute delivery guarantee

## Error Handling

The system handles various error scenarios:
- Insufficient items in inventory
- No available delivery agents
- Transaction timeouts
- Network failures
- Concurrent access conflicts

## Testing

To test different scenarios:

1. Order more items than available:
```go
items := map[string]int{
    "pizza": 10, // More than available
}
```

2. Place multiple orders simultaneously:
```go
for i := 0; i < 3; i++ {
    go func() {
        items := map[string]int{"pizza": 2}
        orderService.PlaceOrder(ctx, items, "location")
    }()
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Two-Phase Commit protocol
- Distributed systems principles
- Food delivery system requirements 