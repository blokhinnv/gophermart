package handlers

// TODO: тесты для tracker + accrualSystem?
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/accrual"
	"github.com/blokhinnv/gophermart/internal/app/database"
	"github.com/blokhinnv/gophermart/internal/app/database/ordertracker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PostOrderTestSuite struct {
	suite.Suite
	AuthHandlerTestSuite
	accrualService *accrual.MockService
	tracker        *ordertracker.MockTracker
}

func (suite *PostOrderTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.db = database.NewMockService(suite.ctrl)
	suite.tracker = ordertracker.NewMockTracker(suite.ctrl)

	suite.accrualService = accrual.NewMockService(suite.ctrl)
	suite.db.EXPECT().
		Tracker().
		Return(suite.tracker)

	postOrder := NewPostOrder(
		suite.db,
		2,
		context.Background(),
		suite.accrualService,
		1*time.Second,
	)
	handler := http.HandlerFunc(postOrder.Handler)
	suite.setupAuth(handler)
}

func (suite *PostOrderTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *PostOrderTestSuite) makeRequest(
	testName string,
	auth bool,
	body io.Reader,
) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/user/orders", body)
	req.Header.Set("Content-Type", "text/plain")
	if auth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", suite.tokenSign))
	}
	suite.handler.ServeHTTP(rr, req)
	log.Printf("[%v]: %v", testName, rr.Body.String())
	return rr
}

func (suite *PostOrderTestSuite) TestNoAuth() {
	rr := suite.makeRequest("TestNoAuth", false, bytes.NewBuffer([]byte(`18`)))
	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *PostOrderTestSuite) TestAccepted() {
	suite.db.EXPECT().
		AddOrder(gomock.Any(), gomock.Eq("18"), gomock.Eq(1)).
		Times(1).
		Return(nil)
	suite.tracker.EXPECT().
		Add(gomock.Any(), gomock.Eq("18")).
		Times(1).
		Return(nil)

	rr := suite.makeRequest("TestAccepted", true, bytes.NewBuffer([]byte(`18`)))
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
	rr := suite.makeRequest("TestAlreadyAddedByMe", true, bytes.NewBuffer([]byte(`18`)))
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *PostOrderTestSuite) TestAlreadyAddedNotByMe() {
	suite.db.EXPECT().
		AddOrder(gomock.Any(), gomock.Eq("18"), gomock.Eq(1)).
		Times(1).
		Return(fmt.Errorf("%w: orderID=%v userID=%v", database.ErrOrderAlreadyAddedByOtherUser, "18", 1))
	rr := suite.makeRequest("TestAlreadyAddedNotByMe", true, bytes.NewBuffer([]byte(`18`)))
	suite.Equal(http.StatusConflict, rr.Code)
}

func (suite *PostOrderTestSuite) TestBadRequest() {
	rr := suite.makeRequest("TestBadRequest", true, bytes.NewBuffer([]byte(``)))
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PostOrderTestSuite) TestBadID() {
	rr := suite.makeRequest("TestBadID", true, bytes.NewBuffer([]byte(`123123123123`)))
	suite.Equal(http.StatusUnprocessableEntity, rr.Code)
}

func TestPostOrderTestSuite(t *testing.T) {
	suite.Run(t, new(PostOrderTestSuite))
}
