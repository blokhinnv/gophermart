package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

func GenerateJWTToken(
	username string,
	signingKey []byte,
	expireDuration time.Duration,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: username,
	})
	return token.SignedString(signingKey)
}
