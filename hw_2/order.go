package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	StatusNew      = "new"
	StatusPaid     = "paid"
	StatusShipped  = "shipped"
	StatusCanceled = "canceled"
)

var (
	ErrInvalidOrderID = errors.New("order id must be positive")
	ErrNilOrder       = errors.New("order is nil")
	ErrEmptyItemName  = errors.New("item name is empty")
	ErrNegativePrice  = errors.New("item price is negative")
)

type OrderID int

type Item struct {
	Name         string
	PriceKopecks int64
}

type Order struct {
	ID     OrderID
	Items  []Item
	Status string
}

func NewOrder(id OrderID) (Order, error) {
	if id <= 0 {
		return Order{}, fmt.Errorf("new order %d: %w", id, ErrInvalidOrderID)
	}
	return Order{ID: id, Status: StatusNew}, nil
}

func AddItem(order *Order, item Item) error {
	if order == nil {
		return fmt.Errorf("add item %q: %w", item.Name, ErrNilOrder)
	}
	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("add item to order %d: %w", order.ID, ErrEmptyItemName)
	}
	if item.PriceKopecks < 0 {
		return fmt.Errorf("add item %q: %w", item.Name, ErrNegativePrice)
	}
	order.Items = append(order.Items, item)
	return nil
}

func Total(order Order) int64 {
	var sum int64
	for _, item := range order.Items {
		sum += item.PriceKopecks
	}
	return sum
}

func Clone(order Order) Order {
	order.Items = slices.Clone(order.Items)
	return order
}

func StatusNames() map[string]string {
	names := make(map[string]string, 4)
	names[StatusNew] = "новый"
	names[StatusPaid] = "оплачен"
	names[StatusShipped] = "отправлен"
	names[StatusCanceled] = "отменен"
	return names
}
