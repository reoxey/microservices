package mock

import (
	"context"

	"product/core"
)

type repo struct {
}

func (r repo) UpdateStocks(ctx context.Context, items core.Items) error {
	return nil
}

func (r repo) All(ctx context.Context) (core.Products, error) {
	return core.Products{
		{
			Id:        1,
			Sku:       "ok1",
			Name:      "One",
			Price:     1.23456,
			CreatedAt: "Today",
		},
		{
			Id:        2,
			Sku:       "ok2",
			Name:      "Two",
			Price:     9.6352,
			CreatedAt: "Today",
		},
	}, nil
}

func (r repo) ByID(ctx context.Context, i int) (*core.Product, error) {
	return &core.Product{
		Id:        i,
		Sku:       "ok1",
		Name:      "One",
		Price:     1.23456,
		CreatedAt: "Today",
	}, nil
}

func (r repo) Add(ctx context.Context, user *core.Product) (int, error) {
	return 0, nil
}

func (r repo) Edit(ctx context.Context, user *core.Product) error {
	return nil
}

func NewMock() core.ProductRepo {
	return &repo{}
}
