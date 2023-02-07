package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/suite"
)

const AccrualSystemAddress = "http://localhost:8081"

type PostOrderTestSuite struct {
	suite.Suite
	db        *database.DatabaseService
	postOrder *PostOrder
	handler   http.Handler
	tokenSign string
}

func (suite *PostOrderTestSuite) SetupSuite() {
	db, _ := database.NewDatabaseService(
		&config.Config{DatabaseURI: "postgres://root:pwd@localhost:5432/root"},
		context.Background(),
		true,
	)
	postOrder := NewPostOrder(db, 10, 2, AccrualSystemAddress)
	suite.db = db
	suite.postOrder = postOrder

	user, pwd := "nikita", "123"
	signingKey := []byte("qwerty")
	suite.db.AddUser(context.Background(), user, pwd)

	token := auth.GenerateJWTToken(
		&models.User{ID: 1, Username: user},
		signingKey,
		time.Hour,
	)

	tokenSign, _ := token.SignedString(signingKey)
	suite.tokenSign = tokenSign
	tokenAuth := jwtauth.New("HS256", signingKey, nil)

	verifier := jwtauth.Verifier(tokenAuth)
	authentifier := jwtauth.Authenticator
	handler := http.HandlerFunc(suite.postOrder.Handler)
	suite.handler = verifier(authentifier(handler))

}

func (suite *PostOrderTestSuite) makeRequest(body io.Reader, auth bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/orders", body)
	req.Header.Set("Content-Type", "text/plain")

	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	return rr
}

func (suite *PostOrderTestSuite) TestNoAuth() {
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`18`)), false)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *PostOrderTestSuite) TestOK() {
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`18`)), true)
	suite.Equal(http.StatusAccepted, rr.Code)
}

func (suite *PostOrderTestSuite) TestBadRequest() {
	rr := suite.makeRequest(bytes.NewBuffer([]byte(``)), true)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PostOrderTestSuite) TestBadID() {
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`123123123123`)), true)
	suite.Equal(http.StatusUnprocessableEntity, rr.Code)
}

func TestPostOrderSuite(t *testing.T) {
	suite.Run(t, new(PostOrderTestSuite))
}
