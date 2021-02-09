package core

type Product struct {
	Id        int     `json:"id"`
	Sku       string  `json:"sku" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Price     float64 `json:"price" validate:"required,gte=1,lte=1000"`
	Stocks    int     `json:"stocks" validate:"required,gt=0"`
	CreatedAt string  `json:"created_at"`
}

type Products []*Product

type Item struct {
	Id  int `json:"id"`
	Qty int `json:"qty"`
}

type Items []*Item
