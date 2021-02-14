package core

import "context"

type CartRepo interface {
	Create(ctx context.Context, buyer int) (int, error)
	ByID(ctx context.Context, cartId int) (Cart, error)
	AddItem(ctx context.Context, cartId int, item *Item) error
	UpdateQty(ctx context.Context, cartId int, item *Item) error
	UpdateItems(ctx context.Context, item *Item) error
	RemoveItem(ctx context.Context, cartId int, itemId int) error
	ResetCart(ctx context.Context, cartId int) error
}
