package security

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

var DB *gorm.DB

type TokenJWT struct{}

func NewTokenJWT() *TokenJWT {
	return &TokenJWT{}
}

type UserClaims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var jwtKey = []byte("a")

func (h *TokenJWT) GenerateTokenJWT(email, password string) (string, error) {
	claims := &UserClaims{
		Email:          email,
		Password:       password,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
