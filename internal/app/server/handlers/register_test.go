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

type RegisterTestSuite struct {
	suite.Suite
	db      *database.DatabaseService
	reg     Register
	handler http.HandlerFunc
}

func (suite *RegisterTestSuite) makeRequest(body io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/register", body)
	req.Header.Set("Content-Type", "application/json")
	suite.handler.ServeHTTP(rr, req)
	return rr

}

func (suite *RegisterTestSuite) SetupSuite() {
	db, _ := database.NewDatabaseService(
		&config.Config{DatabaseURI: "postgres://root:pwd@localhost:5432/root"},
		context.Background(),
		true,
	)
	reg := Register{
		LogReg: LogReg{
			db:             db,
			signingKey:     []byte("qwerty"),
			expireDuration: 1 * time.Hour,
		},
	}
	suite.db = db
	suite.reg = reg
	suite.handler = http.HandlerFunc(suite.reg.Handler)
}

func (suite *RegisterTestSuite) TearDownTest() {
	suite.db.GetConn().Exec(context.Background(), `DELETE FROM UserAccount;`)
}

func (suite *RegisterTestSuite) TestSingle() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	rr := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *RegisterTestSuite) TestAlreadyExisted() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	resp1 := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, resp1.Code)
	resp2 := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusConflict, resp2.Code)
}

func (suite *RegisterTestSuite) TestIncorrentBody() {
	jsonStr := []byte(`{"login":"nikita", "pass`)
	resp1 := suite.makeRequest(bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusBadRequest, resp1.Code)

}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
