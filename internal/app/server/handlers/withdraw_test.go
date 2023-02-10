package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WithdrawTestSuite struct {
	suite.Suite
	AuthHandlerTestSuite
}

func (suite *WithdrawTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	withdraw := Withdraw{db: suite.db}
	handler := http.HandlerFunc(withdraw.Handler)
	suite.setupAuth(handler)
}

func (suite *WithdrawTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *WithdrawTestSuite) makeRequest(
	testName string,
	auth, setContentTypeHeader bool,
	body io.Reader,
) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/balance/withdraw", body)
	if setContentTypeHeader {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *WithdrawTestSuite) TestNoAuth() {
	rr := suite.makeRequest("TestNoAuth", false, false, nil)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *WithdrawTestSuite) TestBadHeader() {
	rr := suite.makeRequest("TestBadHeader", true, false, nil)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *WithdrawTestSuite) TestNotEnoughBalance() {
	jsonStr := []byte(`{"order":"2377225624", "sum": 751}`)

	suite.db.EXPECT().
		GetBalance(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(&models.Balance{
			Current:   sql.NullFloat64{Float64: 100, Valid: true},
			Withdrawn: sql.NullFloat64{Float64: 100, Valid: true},
		}, nil)

	rr := suite.makeRequest("TestNotEnoughBalance", true, true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusPaymentRequired, rr.Code)
}

func (suite *WithdrawTestSuite) TestBadOrderIDFormat() {
	jsonStr := []byte(`{"order":"11111", "sum": 751}`)
	rr := suite.makeRequest("TestBadOrderIDFormat", true, true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusUnprocessableEntity, rr.Code)
}

func (suite *WithdrawTestSuite) TestOK() {
	jsonStr := []byte(`{"order":"18", "sum": 10}`)

	suite.db.EXPECT().
		GetBalance(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(&models.Balance{
			Current:   sql.NullFloat64{Float64: 100, Valid: true},
			Withdrawn: sql.NullFloat64{Float64: 100, Valid: true},
		}, nil)

	suite.db.EXPECT().
		AddWithdrawalRecord(gomock.Any(), gomock.Eq("18"), gomock.Eq(10.0), gomock.Eq(1)).
		Times(1).
		Return(nil)

	rr := suite.makeRequest("TestOK", true, true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, rr.Code)
}

func TestWithdrawTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawTestSuite))
}
