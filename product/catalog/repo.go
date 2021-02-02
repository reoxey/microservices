package catalog

import "context"

type ProductRepo interface {
	All(ctx context.Context) (Products, error)
	ByID(ctx context.Context, id int) (Product, error)
	Add(ctx context.Context, prod *Product) (int, error)
	Edit(ctx context.Context, prod *Product) error
}
