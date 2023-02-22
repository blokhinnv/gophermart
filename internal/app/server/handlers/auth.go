package handlers

import (
	"context"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
)

func GetUserIDFromContext(ctx context.Context) (int, error) {
	_, claims, _ := jwtauth.FromContext(ctx)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("%w: %+v", ErrBadClaims, claims)
	}
	return int(userID), nil
}
