package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WithdrawalsTestSuite struct {
	suite.Suite
	AuthHandlerTestSuite
}

func (suite *WithdrawalsTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	withdrawals := Withdrawals{db: suite.db}
	handler := http.HandlerFunc(withdrawals.Handler)
	suite.setupAuth(handler)
}

func (suite *WithdrawalsTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *WithdrawalsTestSuite) makeRequest(
	testName string,
	auth bool,
) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/widtdrawals", nil)
	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *WithdrawalsTestSuite) TestNoAuth() {
	rr := suite.makeRequest("TestNoAuth", false)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *WithdrawalsTestSuite) TestNoContent() {
	suite.db.EXPECT().
		GetWithdrawals(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(nil, database.ErrEmptyResult)

	rr := suite.makeRequest("TestNoContent", true)
	suite.Equal(http.StatusNoContent, rr.Code)
}

func (suite *WithdrawalsTestSuite) TestMany() {
	start := time.Now()
	withdrawals := []models.Withdrawal{
		{
			Order:       "18",
			Sum:         123,
			ProcessedAt: start,
		},
		{
			Order:       "24",
			Sum:         123,
			ProcessedAt: start.Add(5 * time.Second),
		},
	}

	suite.db.EXPECT().
		GetWithdrawals(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(withdrawals, nil)

	rr := suite.makeRequest("TestMany", true)
	suite.Equal(http.StatusOK, rr.Code)

	expected := fmt.Sprintf(
		`[{"order":"18","sum":123,"processed_at":"%v"},
		{"order":"24","sum":123,"processed_at":"%v"}]`,
		start.Format(time.RFC3339),
		start.Add(5*time.Second).Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func TestWithdrawalsTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawalsTestSuite))
}
