package core

import (
	"context"
	"encoding/json"
)

type orderService struct {
	repo 		OrderRepo
	auth     	JWTService
	publisherQ  Publisher
}

func (o orderService) OrderShipping(ctx context.Context, ship *Shipping) error {

	str, err := json.Marshal(&ship)
	if err != nil {
		return err
	}

	msg := &Message{
		Topic: "shipping",
		Msg: str,
	}

	return o.publisherQ.Publish(ctx, msg)
}

func (o orderService) UpdateItemStocks(ctx context.Context, items []*Item) error {
	byt, err := json.Marshal(&items)
	if err != nil {
		return err
	}
	return o.publisherQ.Publish(ctx, &Message{
		Topic: "product_qty",
		Msg:   byt,
	})
}

func (o orderService) GetOrder(ctx context.Context, orderId int) (*Order, error) {
	return o.repo.GetOrder(ctx, orderId)
}

func (o orderService) AllOrders(ctx context.Context, buyer int) (Orders, error) {
	return o.repo.AllOrders(ctx, buyer)
}

func (o orderService) PlaceOrder(ctx context.Context, order *Order) (int, error) {
	return o.repo.PlaceOrder(ctx, order)
}

func (o orderService) Authorize(s string) (map[string]interface{}, error) {
	return o.auth.ValidateToken(s)
}

func NewService(or OrderRepo, auth JWTService, pub Publisher) OrderService {
	return &orderService{
		or,
		auth,
		pub,
	}
}
