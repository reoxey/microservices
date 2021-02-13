package core

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type shipServ struct {
	repo ShippingRepo
	auth JWTService
	validate *validator.Validate
}

func (s shipServ) AddAddress(ctx context.Context, address *Address) (int, error) {
	if err := s.validate.Struct(address); err != nil {
		return 0, err
	}
	return s.repo.AddAddress(ctx, address)
}

func (s shipServ) AddOrderShipping(ctx context.Context, ship *Shipping) (int, error) {
	addr, err := s.repo.AddressById(ctx, ship.Id)
	if err != nil {
		return 0, err
	}
	ship.Address = addr
	return s.repo.AddOrderShipping(ctx, ship)
}

func (s shipServ) AddressById(ctx context.Context, addressId int) (*Address, error) {
	return s.repo.AddressById(ctx, addressId)
}

func (s shipServ) AllAddresses(ctx context.Context, userId int) (Addresses, error) {
	return s.repo.AllAddresses(ctx, userId)
}

func (s shipServ) EditAddress(ctx context.Context, address *Address) error {
	return s.repo.EditAddress(ctx, address)
}

func (s shipServ) UpdateStatus(ctx context.Context, shipId int, status ShippingStatus) error {
	return s.repo.UpdateStatus(ctx, shipId, status)
}

func (s shipServ) UpdatePayment(ctx context.Context, shipId int, paymentStatus int) error {
	return s.repo.UpdatePayment(ctx, shipId, paymentStatus)
}

func (s shipServ) Authorize(token string) (map[string]interface{}, error) {
	return s.auth.ValidateToken(token)
}

func NewService(sr ShippingRepo, auth JWTService) ShippingService {

	return &shipServ{
		sr,
		auth,
		validator.New(),
	}
}
