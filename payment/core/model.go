package core

type PaymentType int32

const (
	PaymentCard   PaymentType = 0
	PaymentCOD    PaymentType = 1
	PaymentPayPal PaymentType = 2
)

type Payment struct {
	Total float64     `json:"total,omitempty"`
	Type  PaymentType `json:"type"`
}
