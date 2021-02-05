package mock

import (
	"context"

	"shipping/core"
)

type repo struct {
}

func (r repo) AddAddress(ctx context.Context, address *core.Address) (int, error) {
	return 0, nil
}

func (r repo) AddressById(ctx context.Context, id int) (*core.Address, error) {
	return &core.Address{
			Id:           0,
			User:         0,
			ContactName:  "",
			ContactPhone: "",
			Landmark:     "",
			City:         "",
			State:        "",
			Country:      "",
			Zip:          0,
			CreatedAt:    "",
		},
		nil
}

func (r repo) AllAddresses(ctx context.Context) (core.Addresses, error) {
	return core.Addresses{}, nil
}

func (r repo) EditAddress(ctx context.Context, address *core.Address) error {
	return nil
}

func (r repo) UpdateStatus(ctx context.Context, id int, status core.ShippingStatus) error {
	return nil
}

func (r repo) UpdatePayment(ctx context.Context, id int, paymentStatus int) error {
	return nil
}

func NewMock() core.ShippingRepo {
	return &repo{}
}
