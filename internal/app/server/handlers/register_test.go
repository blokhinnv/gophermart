package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

type RegisterTestSuite struct {
	suite.Suite
	db      *database.MockService
	handler http.HandlerFunc
	ctrl    *gomock.Controller
}

func (suite *RegisterTestSuite) makeRequest(
	testName string,
	setContentTypeHeader bool,
	body io.Reader,
) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/register", body)
	if setContentTypeHeader {
		req.Header.Set("Content-Type", "application/json")
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *RegisterTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)

	reg := Register{
		LogReg: LogReg{
			db:             suite.db,
			signingKey:     []byte("qwerty"),
			expireDuration: 1 * time.Hour,
		},
	}
	suite.handler = http.HandlerFunc(reg.Handler)
}

func (suite *RegisterTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *RegisterTestSuite) TestSingle() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	suite.db.EXPECT().
		AddUser(gomock.Any(), gomock.Eq("nikita"), gomock.Eq("123")).
		Times(1).
		Return(&models.User{
			ID:             1,
			Username:       "nikita",
			HashedPassword: auth.GenerateHash("123", "456"),
			Salt:           "456",
		}, nil)

	rr := suite.makeRequest("TestSingle", true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *RegisterTestSuite) TestAlreadyExisted() {
	jsonStr := []byte(`{"login":"nikita", "password": "123"}`)
	suite.db.EXPECT().
		AddUser(gomock.Any(), gomock.Eq("nikita"), gomock.Eq("123")).
		Times(1).
		Return(&models.User{
			ID:             1,
			Username:       "nikita",
			HashedPassword: auth.GenerateHash("123", "456"),
			Salt:           "456",
		}, nil)
	resp1 := suite.makeRequest("TestAlreadyExisted", true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusOK, resp1.Code)
	suite.db.EXPECT().
		AddUser(gomock.Any(), gomock.Eq("nikita"), gomock.Eq("123")).
		Times(1).
		Return(nil, fmt.Errorf("%w: %v", database.ErrUserAlreadyExists, "nikita"))
	resp2 := suite.makeRequest("TestAlreadyExisted", true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusConflict, resp2.Code)
}

func (suite *RegisterTestSuite) TestIncorrentBody() {
	jsonStr := []byte(`{"login":"nikita", "pass`)
	resp1 := suite.makeRequest("TestIncorrentBody", true, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusBadRequest, resp1.Code)
}

func (suite *RegisterTestSuite) TestNoContentType() {
	jsonStr := []byte(`{"login":"nikita", "pass`)
	resp1 := suite.makeRequest("TestNoContentType", false, bytes.NewBuffer(jsonStr))
	suite.Equal(http.StatusBadRequest, resp1.Code)

}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
