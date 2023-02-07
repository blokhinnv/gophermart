package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/auth"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	db      *database.MockService
	handler http.HandlerFunc
	ctrl    *gomock.Controller
}

func (suite *LoginTestSuite) makeRequest(body io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/login", body)
	req.Header.Set("Content-Type", "application/json")
	suite.handler.ServeHTTP(rr, req)
	fmt.Println(rr.Body.String())
	return rr

}

func (suite *LoginTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	log := Login{
		LogReg: LogReg{
			db:             suite.db,
			signingKey:     []byte("qwerty"),
			expireDuration: 1 * time.Hour,
		},
	}
	suite.handler = http.HandlerFunc(log.Handler)
}

func (suite *LoginTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *LoginTestSuite) TestOk() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	suite.db.EXPECT().
		FindUser(gomock.Any(), gomock.Eq("nikita"), gomock.Eq("123")).
		Times(1).
		Return(&models.User{
			ID:             1,
			Username:       "nikita",
			HashedPassword: auth.GenerateHash("123", "456"),
			Salt:           "456",
		}, nil)

	rr := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *LoginTestSuite) TestWrong() {
	jsonStr := []byte(`{"login":"nikita", "password": "1234"}`)
	suite.db.EXPECT().
		FindUser(gomock.Any(), gomock.Eq("nikita"), gomock.Eq("1234")).
		Times(1).
		Return(&models.User{
			ID:             1,
			Username:       "nikita",
			HashedPassword: auth.GenerateHash("123", "456"),
			Salt:           "456",
		}, nil)

	resp1 := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusUnauthorized, resp1.Code)
}

func (suite *LoginTestSuite) TestIncorrentBody() {
	jsonStr := []byte(`{"login":"nikita", "pass`)
	resp1 := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusBadRequest, resp1.Code)

}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
