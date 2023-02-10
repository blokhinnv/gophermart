package auth

import (
	"time"

	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(
	user *models.User,
	signingKey []byte,
	expireDuration time.Duration,
) *jwt.Token {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   user.ID,
		Username: user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}
