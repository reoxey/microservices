package core

type payService struct {
	auth JWTService
}

func (p payService) Authorize(token string) (map[string]interface{}, error) {
	panic("implement me")
}

func NewService(auth JWTService) PayService {
	return &payService{
		auth,
	}
}
