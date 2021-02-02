package mock

import (
	"context"

	"product/catalog"
)

type repo struct {

}

func (r repo) All(ctx context.Context) (catalog.Products, error) {
	return []catalog.Product{
		{
			Id:       1,
			Sku: "ok1",
			Name:     "One",
			Price: 1.23456,
			CreatedAt: "Today",
		},
		{
			Id:       2,
			Sku: "ok2",
			Name:     "Two",
			Price: 9.6352,
			CreatedAt: "Today",
		},
	}, nil
}

func (r repo) ByID(ctx context.Context, i int) (catalog.Product, error) {
	return catalog.Product{
		Id:       i,
		Sku: "ok1",
		Name:     "One",
		Price: 1.23456,
		CreatedAt: "Today",
	}, nil
}

func (r repo) Add(ctx context.Context, user *catalog.Product) (int, error) {
	return 0, nil
}

func (r repo) Edit(ctx context.Context, user *catalog.Product) error {
	return nil
}

func NewMock() catalog.ProductRepo {
	return &repo{}
}
