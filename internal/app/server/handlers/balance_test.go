package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BalanceTestSuite struct {
	AuthHandlerTestSuite
}

func (suite *BalanceTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	balance := Balance{db: suite.db}
	handler := http.HandlerFunc(balance.Handler)
	suite.setupAuth(handler)
}

func (suite *BalanceTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *BalanceTestSuite) makeRequest(testName string, auth bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req.Header.Set("Content-Type", "text/plain")
	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *BalanceTestSuite) TestNoAuth() {
	rr := suite.makeRequest("TestNoAuth", false)
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

	rr := suite.makeRequest("TestNoTransactions", true)
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

	rr := suite.makeRequest("TestSomeTransactions", true)
	suite.Equal(http.StatusOK, rr.Code)
	expected := `{
		"current": 250.32,
		"withdrawn": 50.54
	}`
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceTestSuite))
}
