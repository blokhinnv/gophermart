package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PostOrderTestSuite struct {
	suite.Suite
	db             *database.MockService
	handler        http.Handler
	tokenSign      string
	ctrl           *gomock.Controller
	accrualService *accrual.MockService
}

func (suite *PostOrderTestSuite) SetupSuite() {
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

	suite.accrualService = accrual.NewMockService(suite.ctrl)
	suite.accrualService.EXPECT().
		GetOrderInfo(gomock.Eq("18")).
		Return([]byte(`{
			"order": "18",
			"status": "PROCESSED",
			"accrual": 500
		}`), nil)

	postOrder := NewPostOrder(suite.db, 10, 2, suite.accrualService)
	verifier := jwtauth.Verifier(tokenAuth)
	authentifier := jwtauth.Authenticator
	handler := http.HandlerFunc(postOrder.Handler)
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

func (suite *PostOrderTestSuite) TestAccepted() {
	suite.db.EXPECT().
		AddOrder(gomock.Any(), gomock.Eq("18"), gomock.Eq(1)).
		Times(1).
		Return(nil)
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`18`)), true)
	suite.Equal(http.StatusAccepted, rr.Code)
}

func (suite *PostOrderTestSuite) TestAlreadyAddedByMe() {
	suite.db.EXPECT().
		AddOrder(gomock.Any(), gomock.Eq("18"), gomock.Eq(1)).
		Times(1).
		Return(fmt.Errorf(
			"%w: orderID=%v userID=%v",
			database.ErrOrderAlreadyAddedByThisUser,
			"18",
			1,
		))
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`18`)), true)
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *PostOrderTestSuite) TestAlreadyAddedNotByMe() {
	suite.db.EXPECT().
		AddOrder(gomock.Any(), gomock.Eq("18"), gomock.Eq(1)).
		Times(1).
		Return(fmt.Errorf("%w: orderID=%v userID=%v", database.ErrOrderAlreadyAddedByOtherUser, "18", 1))
	rr := suite.makeRequest(bytes.NewBuffer([]byte(`18`)), true)
	suite.Equal(http.StatusConflict, rr.Code)
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
