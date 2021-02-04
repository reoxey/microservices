package jwtauth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"user/core"
)

type authClaims struct {
	Id int `json:"id"`
	Email string `json:"email"`
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issure string
}

func (j jwtServices) GenerateToken(id int, email string, isAdmin bool) (string, error) {
	claims := &authClaims{
		id,
		email,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    j.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j jwtServices) ValidateToken(token string) (map[string]interface{}, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])

		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, err
	}
	return jwtToken.Claims.(jwt.MapClaims), nil
}

func New() core.JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure: "reoxey",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}
