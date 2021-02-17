package mock

import (
	"context"

	"order/core"
)

type repo struct {

}

func (r repo) GetOrder(ctx context.Context, orderId int) (*core.Order, error) {
	return &core.Order{}, nil
}

func (r repo) AllOrders(ctx context.Context, buyer int) (core.Orders, error) {
	return core.Orders{}, nil
}

func (r repo) PlaceOrder(ctx context.Context, order *core.Order) (id int, err error) {
	return
}

func NewMock() core.OrderRepo {
	return &repo{}
}
