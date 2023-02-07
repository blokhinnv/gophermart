package auth

import (
	"time"

	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/dgrijalva/jwt-go/v4"
)

func GenerateJWTToken(
	user *models.User,
	signingKey []byte,
	expireDuration time.Duration,
) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserID:   user.ID,
		Username: user.Username,
	})
	return token
}
