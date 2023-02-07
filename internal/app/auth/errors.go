package auth

import "errors"

var ErrInvalidAccessToken = errors.New("incorrent jwt token")
