package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"

	"cart/catalogpb"
)

type cartService struct {
	repo     CartRepo
	auth     JWTService
	validate *validator.Validate
	grpcCatalog catalogpb.CatalogClient
	publisherQ Publisher
}

func (c cartService) Authorize(s string) (map[string]interface{}, error) {
	return c.auth.ValidateToken(s)
}

func (c cartService) New(ctx context.Context, buyer int) (int, error) {
	return c.repo.Create(ctx, buyer)
}

func (c cartService) AddToCart(ctx context.Context, cartId int, item *Item) (err error) {
	if err = c.validate.Struct(item); err != nil {
		return
	}

	prod, err := c.grpcCatalog.GetProduct(ctx,
		&catalogpb.ProductId{Id: int32(item.Id)},
		)
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", prod)

	if int(prod.Stocks) < item.Qty {
		return fmt.Errorf("not enough stocks for %d; qty %d/%d stocks", item.Id, item.Qty, prod.Stocks)
	}

	item.Name = prod.Name
	item.Price = prod.Price
	item.Stocks = int(prod.Stocks)

	if err = c.repo.AddItem(ctx, cartId, item); err != nil {
		return
	}
	return
}

func (c cartService) Show(ctx context.Context, cartId int) (Cart, error) {
	cart, err := c.repo.ByID(ctx, cartId)
	if err != nil {
		return Cart{}, err
	}
	sum := 0.0
	for _, item := range cart.Items {
		sum += item.Price * float64(item.Qty)
	}
	cart.Payment = &Payment{sum}

	return cart, nil
}

func (c cartService) UpdateQty(ctx context.Context, cartId int, item *Item) error {
	return c.repo.UpdateQty(ctx, cartId, item)
}

func (c cartService) UpdateItems(ctx context.Context, item *Item) error {
	return c.repo.UpdateItems(ctx, item)
}

func (c cartService) DeleteItems(ctx context.Context, cartId int, itemId int) error {
	return c.repo.RemoveItem(ctx, cartId, itemId)
}

func (c cartService) Checkout(ctx context.Context, checkout *Checkout, cartId int) error {

	cart, err := c.Show(ctx, cartId)
	if err != nil {
		return err
	}

	order := &Order{
		Cart:     &cart,
		Checkout: checkout,
	}

	str, err := json.Marshal(&order)
	if err != nil {
		return err
	}

	if err = c.publisherQ.Publish(ctx, &Message{
		Topic: "order",
		Msg:   str,
	}); err != nil {
		return err
	}

	return c.repo.ResetCart(ctx, cartId)
}

func NewService(cr CartRepo, auth JWTService, gClient catalogpb.CatalogClient, publisher Publisher) CartService {
	return &cartService{
		cr,
		auth,
		validator.New(),
		gClient,
		publisher,
	}
}
