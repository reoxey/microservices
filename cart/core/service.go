package core

import "context"

type CartService interface {
	New(ctx context.Context, buyer int) (int, error)
	AddToCart(ctx context.Context, cartId int, item *Item) error
	Show(ctx context.Context, cartId int) (Cart, error)
	UpdateQty(ctx context.Context, cartId int, item *Item) error
	UpdateItems(ctx context.Context, item *Item) error
	DeleteItems(ctx context.Context, cartId int, itemId int) error
	Checkout(ctx context.Context, checkout *Checkout, cartId int) error
	Authorize(token string) (map[string]interface{}, error)
}
