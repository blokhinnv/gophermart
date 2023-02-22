package handlers

import "errors"

var ErrIncorrectRequest = errors.New("incorrent request")
var ErrIncorrectContentType = errors.New("incorrent content-type")
var ErrNotValid = errors.New("data is not valid")
var ErrIncorrectCredentials = errors.New("incorrent credentials")
var ErrNotEnoughBalance = errors.New("not enough points on balance")
var ErrBadClaims = errors.New("incorrect claims")
var ErrServerShutdown = errors.New("server is shutting down")
