package core

import "context"

type ShippingService interface {
	AddAddress(ctx context.Context, address *Address) (int, error)
	AddOrderShipping(ctx context.Context, ship *Shipping) (int, error)
	AddressById(ctx context.Context, addressId int) (*Address, error)
	AllAddresses(ctx context.Context, userId int) (Addresses, error)
	EditAddress(ctx context.Context, address *Address) error
	UpdateStatus(ctx context.Context, shipId int, status ShippingStatus) error
	UpdatePayment(ctx context.Context, shipId int, paymentStatus int) error
	Authorize(token string) (map[string]interface{}, error)
}

