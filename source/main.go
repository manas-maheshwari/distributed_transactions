package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"distributed_transactions/source/coordinator"
	"distributed_transactions/source/delivery"
	"distributed_transactions/source/orders"
	"distributed_transactions/source/store"
)

func main() {
	// Initialize services
	coord := coordinator.NewCoordinator()
	storeService := store.NewStore()
	deliveryService := delivery.NewDelivery()
	orderService := orders.NewOrderService(coord, storeService, deliveryService)

	// Add some initial data
	storeService.AddItem("pizza", 5)
	storeService.AddItem("burger", 10)
	deliveryService.AddAgent("agent1", "location1")
	deliveryService.AddAgent("agent2", "location2")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Example: Place an order
	items := map[string]int{
		"pizza": 2,
		"burger": 1,
	}

	order, err := orderService.PlaceOrder(ctx, items, "customer_location")
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	fmt.Printf("Order placed successfully! Order ID: %s\n", order.ID)
	fmt.Printf("Order status: %s\n", order.Status)
	fmt.Printf("Order items: %v\n", order.Items)
	fmt.Printf("Delivery location: %s\n", order.Location)

	// Example: Get order status
	retrievedOrder, err := orderService.GetOrder(ctx, order.ID)
	if err != nil {
		log.Fatalf("Failed to get order: %v", err)
	}

	fmt.Printf("\nRetrieved order status: %s\n", retrievedOrder.Status)
} 