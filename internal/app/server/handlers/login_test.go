package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/server/config"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	db      *database.DatabaseService
	log     Login
	handler http.HandlerFunc
}

func (suite *LoginTestSuite) makeRequest(body io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/login", body)
	req.Header.Set("Content-Type", "application/json")
	suite.handler.ServeHTTP(rr, req)
	return rr

}

func (suite *LoginTestSuite) SetupSuite() {
	db, _ := database.NewDatabaseService(
		&config.Config{DatabaseURI: "postgres://root:pwd@localhost:5432/root"},
		context.Background(),
		true,
	)
	reg := Login{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte("qwerty"),
			expireDuration: 1 * time.Hour,
		},
	}
	suite.db = db
	suite.log = reg
	suite.handler = http.HandlerFunc(suite.log.Handler)
	suite.db.AddUser(context.Background(), "nikita", "123")
}

func (suite *LoginTestSuite) TestOk() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	rr := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *LoginTestSuite) TestWrong() {
	jsonStr := []byte(`{"login":"nikita", "password": "1234"}`)
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
