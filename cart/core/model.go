package core

type Cart struct {
	Id        int      `json:"id" validate:"required"`
	Buyer     int      `json:"buyer" validate:"required"`
	Items     []*Item  `json:"items"`
	Payment   *Payment `json:"payment"`
	CreatedAt string   `json:"created_at"`
}

type Item struct {
	Id       int     `json:"id" validate:"required,gt=0"`
	Name     string  `json:"name"`
	OldPrice float64 `json:"old_price,omitempty"`
	Price    float64 `json:"price"`
	Stocks   int     `json:"stocks,omitempty"`
	Qty      int     `json:"qty" validate:"required,gt=0"`
}

type Payment struct {
	Total float64 `json:"total"`
}

type PaymentType int
const (
	Card PaymentType = iota
	COD
	PayPal
)

type Checkout struct {
	AddressId int `json:"address_id"`
	PaymentMethod PaymentType `json:"payment_method"`
}

type Order struct {
	*Cart
	*Checkout
}
