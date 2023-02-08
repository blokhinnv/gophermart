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

type GetOrdersTestSuite struct {
	suite.Suite
	db        *database.MockService
	handler   http.Handler
	tokenSign string
	ctrl      *gomock.Controller
}

func (suite *GetOrdersTestSuite) SetupSuite() {
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

	postOrder := GetOrder{db: suite.db}
	verifier := jwtauth.Verifier(tokenAuth)
	authentifier := jwtauth.Authenticator
	handler := http.HandlerFunc(postOrder.Handler)
	suite.handler = verifier(authentifier(handler))

}

func (suite *GetOrdersTestSuite) makeRequest(auth bool) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req.Header.Set("Content-Type", "text/plain")

	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	return rr
}

func (suite *GetOrdersTestSuite) TestNoAuth() {
	rr := suite.makeRequest(false)
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *GetOrdersTestSuite) TestNoContent() {
	suite.db.EXPECT().
		FindOrdersByUserID(gomock.Any(), gomock.Eq(1)).
		Times(1).
		Return(nil, database.ErrEmptyResult)

	rr := suite.makeRequest(true)
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

	rr := suite.makeRequest(true)

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

	rr := suite.makeRequest(true)

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

	rr := suite.makeRequest(true)

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

	rr := suite.makeRequest(true)

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

	rr := suite.makeRequest(true)

	suite.Equal(http.StatusOK, rr.Code)
	expected := fmt.Sprintf(
		`[{"number":"18","status":"PROCESSED","uploaded_at":"%v", "accrual": 81.4},
		{"number":"24","status":"NEW","uploaded_at":"%v"}]`,
		start.Format(time.RFC3339),
		start.Add(5*time.Second).Format(time.RFC3339),
	)
	assert.JSONEq(suite.T(), expected, rr.Body.String())
}

func TestGetOrdersSuite(t *testing.T) {
	suite.Run(t, new(GetOrdersTestSuite))
}
