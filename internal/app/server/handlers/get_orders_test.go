package handlers

import (
	"database/sql"
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

type GetOrdersTestSuite struct {
	AuthHandlerTestSuite
}

func (suite *GetOrdersTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	getOrder := GetOrder{db: suite.db}
	handler := http.HandlerFunc(getOrder.Handler)
	suite.setupAuth(handler)
}

func (suite *GetOrdersTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *GetOrdersTestSuite) makeRequest(
	testName string,
	auth bool,
) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Content-Type", "text/plain")
	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *GetOrdersTestSuite) TestNoAuth() {
	rr := suite.makeRequest("TestNoAuth", false)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *GetOrdersTestSuite) TestNoContent() {
	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(nil, database.ErrEmptyResult)

	rr := suite.makeRequest("TestNoContent", true)
	suite.Equal(http.StatusNoContent, rr.Code)
}

func (suite *GetOrdersTestSuite) TestNew() {
	start := time.Now()
	orders := []models.Order{
		{
			ID:         "18",
			Status:     "NEW",
			UploadedAt: start,
			Accrual:    sql.NullFloat64{Float64: 0, Valid: false},
		},
	}

	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(orders, nil)

	rr := suite.makeRequest("TestNew", true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"NEW","uploaded_at":"%v"}]`,
		start.Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func (suite *GetOrdersTestSuite) TestProcessing() {
	start := time.Now()
	orders := []models.Order{
		{
			ID:         "18",
			Status:     "PROCESSING",
			UploadedAt: start,
			Accrual:    sql.NullFloat64{Float64: 0, Valid: false},
		},
	}

	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(orders, nil)

	rr := suite.makeRequest("TestProcessing", true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"PROCESSING","uploaded_at":"%v"}]`,
		start.Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func (suite *GetOrdersTestSuite) TestInvalid() {
	start := time.Now()
	orders := []models.Order{
		{
			ID:         "18",
			Status:     "INVALID",
			UploadedAt: start,
			Accrual:    sql.NullFloat64{Float64: 0, Valid: false},
		},
	}

	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(orders, nil)

	rr := suite.makeRequest("TestInvalid", true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"INVALID","uploaded_at":"%v"}]`,
		start.Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func (suite *GetOrdersTestSuite) TestProcessed() {
	start := time.Now()
	orders := []models.Order{
		{
			ID:         "18",
			Status:     "PROCESSED",
			UploadedAt: start,
			Accrual:    sql.NullFloat64{Float64: 500.5, Valid: true},
		},
	}

	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(orders, nil)

	rr := suite.makeRequest("TestProcessed", true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"PROCESSED","uploaded_at":"%v", "accrual": 500.5}]`,
		start.Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func (suite *GetOrdersTestSuite) TestMany() {
	start := time.Now()
	orders := []models.Order{
		{
			ID:         "18",
			Status:     "PROCESSED",
			UploadedAt: start,
			Accrual:    sql.NullFloat64{Float64: 81.4, Valid: true},
		},
		{
			ID:         "24",
			Status:     "NEW",
			UploadedAt: start.Add(5 * time.Second),
			Accrual:    sql.NullFloat64{Float64: 0, Valid: false},
		},
	}

	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(orders, nil)

	rr := suite.makeRequest("TestMany", true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"PROCESSED","uploaded_at":"%v", "accrual": 81.4},
		{"number":"24","status":"NEW","uploaded_at":"%v"}]`,
		start.Format(time.RFC3339),
		start.Add(5*time.Second).Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func TestGetOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersTestSuite))
}
