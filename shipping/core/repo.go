package core

import "context"

type ShippingRepo interface {
	AddAddress(ctx context.Context, address *Address) (int, error)
	AddressById(ctx context.Context, addressId int) (*Address, error)
	AllAddresses(ctx context.Context) (Addresses, error)
	EditAddress(ctx context.Context, address *Address) error
	UpdateStatus(ctx context.Context, shipId int, status ShippingStatus) error
	UpdatePayment(ctx context.Context, shipId int, paymentStatus int) error
}
