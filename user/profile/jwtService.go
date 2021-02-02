package profile

type JWTService interface {
	GenerateToken(id int, email string, isAdmin bool) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}
