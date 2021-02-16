package core

import "context"

type OrderRepo interface {
	GetOrder(ctx context.Context, orderId int) (*Order, error)
	AllOrders(ctx context.Context, buyer int) (Orders, error)
	PlaceOrder(ctx context.Context, order *Order) (int, error)
}
