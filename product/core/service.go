package core

import (
	"context"
)

type ProductService interface {
	AllProducts(ctx context.Context) (Products, error)
	ProductById(ctx context.Context, id int) (*Product, error)
	AddProduct(ctx context.Context, prod *Product) (int, error)
	EditProduct(ctx context.Context, prod *Product) error
	UpdateProductStocks(ctx context.Context, items Items) error
	Authorize(string) (map[string]interface{}, error)
}
