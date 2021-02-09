package catalogpb

import (
	"context"
	"fmt"

	"product/core"
)

type Server struct {
	service core.ProductService
}

func NewGrpc(ps core.ProductService) *Server {
	return &Server{ps}
}

func (s *Server) GetProduct(ctx context.Context, in *ProductId) (*Product, error) {
	prod, err := s.service.ProductById(ctx, int(in.Id))
	if err != nil {
		return nil, err
	}

	if prod == nil {
		return nil, fmt.Errorf("invalid_prod %d", in.Id)
	}

	fmt.Println(prod)

	p := &Product{
		Id:     int32(prod.Id),
		Name:   prod.Name,
		Price:  prod.Price,
		Stocks: int32(prod.Stocks),
	}

	return p, err
}
