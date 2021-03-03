package core

type PayService interface {
	Authorize(token string) (map[string]interface{}, error)
}
