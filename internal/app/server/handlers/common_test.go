package handlers

import (
	"net/http"
	"time"

	_ "github.com/blokhinnv/gophermart/internal/app"
	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	db        *database.MockService
	handler   http.Handler
	tokenSign string
	ctrl      *gomock.Controller
}

func (suite *AuthHandlerTestSuite) setupAuth(
	handler http.HandlerFunc,
) {
	signingKey := []byte("qwerty")
	token := auth.GenerateJWTToken(
		&models.User{ID: 1, Username: "nikita"},
		signingKey,
		time.Hour,
	)
	tokenSign, _ := token.SignedString(signingKey)
	suite.tokenSign = tokenSign
	tokenAuth := jwtauth.New("HS256", signingKey, nil)
	verifier := jwtauth.Verifier(tokenAuth)
	authentifier := jwtauth.Authenticator
	suite.handler = verifier(authentifier(handler))
}
