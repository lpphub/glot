package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const JwtSecret = "wxx87x"

type dataClaims struct {
	jwt.RegisteredClaims
	Data string
}

func GenerateToken(data, secret string) (string, error) {
	return GenerateTokenWithExpire(data, secret, time.Now().Add(2*time.Hour))
}

func GenerateTokenWithExpire(data, secret string, expireIn time.Time) (string, error) {
	c := dataClaims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireIn),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tokenClaims.SignedString([]byte(secret))
}

func ParseToken(token, secret string) (string, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &dataClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if c, ok := tokenClaims.Claims.(*dataClaims); ok && tokenClaims.Valid {
		return c.Data, nil
	}
	return "", err
}
