package handlers

import "errors"

var ErrIncorrectRequest = errors.New("incorrent request")
var ErrIncorrectCredentials = errors.New("incorrent credentials")
