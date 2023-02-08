package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTestSuite struct {
	suite.Suite
	db        *database.MockService
	handler   http.Handler
	tokenSign string
	ctrl      *gomock.Controller
}

func (suite *BalanceTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)

	signingKey := []byte("qwerty")
	token := auth.GenerateJWTToken(
		&models.User{ID: 1, Username: "nikita"},
		signingKey,
		time.Hour,
	)
	tokenSign, _ := token.SignedString(signingKey)
	suite.tokenSign = tokenSign
	tokenAuth := jwtauth.New("HS256", signingKey, nil)

	postOrder := Balance{db: suite.db}
	verifier := jwtauth.Verifier(tokenAuth)
	authentifier := jwtauth.Authenticator
	handler := http.HandlerFunc(postOrder.Handler)
	suite.handler = verifier(authentifier(handler))

}

func (suite *BalanceTestSuite) makeRequest(auth bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req.Header.Set("Content-Type", "text/plain")

	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	return rr
}

func (suite *BalanceTestSuite) TestNoAuth() {
	rr := suite.makeRequest(false)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *BalanceTestSuite) TestNoTransactions() {
	suite.db.EXPECT().
		GetBalance(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(&models.Balance{
			Current:   sql.NullFloat64{Float64: 0, Valid: false},
			Withdrawn: sql.NullFloat64{Float64: 0, Valid: false},
		}, nil)

	rr := suite.makeRequest(true)
	suite.Equal(http.StatusOK, rr.Code)
	expected := `{
		"current": 0,
		"withdrawn": 0
	}`
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func (suite *BalanceTestSuite) TestSomeTransactions() {
	suite.db.EXPECT().
		GetBalance(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(&models.Balance{
			Current:   sql.NullFloat64{Float64: 250.32, Valid: true},
			Withdrawn: sql.NullFloat64{Float64: 50.54, Valid: true},
		}, nil)

	rr := suite.makeRequest(true)
	suite.Equal(http.StatusOK, rr.Code)
	expected := `{
		"current": 250.32,
		"withdrawn": 50.54
	}`
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func TestBalanceSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}
