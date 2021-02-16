package core

type OrderStatus int32

const (
	OrderCreated   OrderStatus = 0
	OrderPaid      OrderStatus = 1
	OrderCanceled  OrderStatus = 2
	OrderFulfilled OrderStatus = 3
	OrderReturned  OrderStatus = 4
)

type PaymentType int32

const (
	PaymentCard   PaymentType = 0
	PaymentCOD    PaymentType = 1
	PaymentPayPal PaymentType = 2
)

type Order struct {
	Id        int         `json:"id,omitempty"`
	CartId    int         `json:"cart_id"`
	Buyer     int         `json:"buyer,omitempty"`
	Items     []*Item     `json:"items,omitempty"`
	Payment   *Payment    `json:"payment,omitempty"`
	Status    OrderStatus `json:"status,omitempty"`
	CreatedAt string      `json:"created_at"`
}

type Payment struct {
	Total float64     `json:"total,omitempty"`
	Type  PaymentType `json:"type,omitempty"`
}

type Item struct {
	Id    int     `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`
	Qty   int32   `json:"qty,omitempty"`
}

//type Cart struct {
//	Buyer     int      `json:"buyer" validate:"required"`
//	Items     []*Item  `json:"items"`
//	Payment   *Payment `json:"payment"`
//	CreatedAt string   `json:"created_at"`
//}

type Checkout struct {
	AddressId     int         `json:"address_id"`
	PaymentMethod PaymentType `json:"payment_method"`
}

type CartOrder struct {
	*Order
	*Checkout
}

type ShippingStatus int

const (
	Ordered ShippingStatus = iota
	Dispatched
	InTransit
	Delivered
	RequestedReturn
	Returned
)

type Shipping struct {
	OrderId       int            `json:"order_id"`
	Status        ShippingStatus `json:"status"`
	PaymentHandle int            `json:"payment_handle"`
	PaymentStatus int            `json:"payment_status"`
	*Address
}

type Address struct {
	Id int `json:"id"`
}

type Orders []*Order
