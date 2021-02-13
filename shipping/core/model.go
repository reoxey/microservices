package core

type ShippingStatus int

const (
	Ordered ShippingStatus = iota
	Dispatched
	InTransit
	Delivered
	RequestedReturn
	Returned
)

type Address struct {
	Id           int    `json:"id"`
	User		 int    `json:"user" validate:"required"`
	ContactName  string `json:"contact_name" validate:"required"`
	ContactPhone string `json:"contact_phone" validate:"required"`
	Landmark     string `json:"landmark" validate:"required"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	Country      string `json:"country" validate:"required"`
	Zip          int    `json:"zip" validate:"required"`
	CreatedAt 	 string `json:"created_at"`
}

type Shipping struct {
	Id            int            `json:"id"`
	OrderId       int            `json:"order_id"`
	Status        ShippingStatus `json:"status"`
	PaymentHandle int            `json:"payment_handle"`
	PaymentStatus int            `json:"payment_status"`
	*Address
}

type Addresses []*Address
