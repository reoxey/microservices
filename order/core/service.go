package core

import "context"

type OrderService interface {
	GetOrder(ctx context.Context, orderId int) (*Order, error)
	AllOrders(ctx context.Context, buyer int) (Orders, error)
	PlaceOrder(ctx context.Context, order *Order) (int, error)
	OrderShipping(ctx context.Context, ship *Shipping) error
	UpdateItemStocks(ctx context.Context, items []*Item) error
	Authorize(token string) (map[string]interface{}, error)
}
