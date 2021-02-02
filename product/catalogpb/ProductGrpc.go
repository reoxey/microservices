package catalogpb

import (
	"context"
	"fmt"

	"product/catalog"
)

type Server struct {
	service catalog.ProductService
}

func NewGrpc(ps catalog.ProductService) *Server {
	return &Server{ps}
}

func (s *Server) GetProduct(ctx context.Context, in *ProductId) (*Product, error) {
	prod, err := s.service.ProductById(ctx, int(in.Id))
	if err != nil {
		return nil, err
	}

	fmt.Println(prod)

	p := &Product{
		Id:            int32(prod.Id),
		Name:          prod.Name,
		Price:         prod.Price,
		Stocks:        int32(prod.Stocks),
	}

	return p, err
}
