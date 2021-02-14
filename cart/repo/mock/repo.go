package mock

import (
	"context"

	"cart/core"
)

type repo struct {

}

func (r repo) UpdateQty(ctx context.Context, cartId int, item *core.Item) error {
	return nil
}

func (r repo) UpdateItems(ctx context.Context, item *core.Item) error {
	return nil
}

func (r repo) Create(ctx context.Context, i int) (int, error) {
	return 0, nil
}

func (r repo) ByID(ctx context.Context, i int) (core.Cart, error) {

	return core.Cart{
		Items:     []*core.Item{
			{

			},
		},
	}, nil
}

func (r repo) AddItem(ctx context.Context, i int, item *core.Item) error {
	return nil
}

func (r repo) RemoveItem(ctx context.Context, i int, i2 int) error {
	return nil
}

func (r repo) ResetCart(ctx context.Context, cartId int) error {
	return nil
}

func NewMock() core.CartRepo {
	return &repo{}
}
