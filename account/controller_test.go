package account_test

import (
	"bank-api/account"
	mock_account "bank-api/account/mocks"
	"bank-api/config"
	"bank-api/helper"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ControllerTestSuite struct {
	suite.Suite
	accountController  *account.AccountController
	mockAccountService *mock_account.MockIAccountService
	mockCtrl           *gomock.Controller
	router             *gin.Engine
}

func (suite *ControllerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockAccountService = mock_account.NewMockIAccountService(ctrl)
	suite.router = config.InitRouter()
	suite.accountController = account.NewAccountController(suite.mockAccountService)
	suite.setEndpoints()
}

func (suite *ControllerTestSuite) setEndpoints() {
	v1 := suite.router.Group("")
	{
		v1.POST(account.AccountEndpoints.NewAccount, suite.accountController.NewAccount)
		v1.GET(account.AccountEndpoints.FindByID, suite.accountController.FindByID)
	}
}

func TestExampleSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (suite *ControllerTestSuite) TestNewAccount() {
	objectID := primitive.NewObjectID()
	data := []struct {
		accountCreated   *account.Account
		url              string
		method           string
		bodyData         string
		expectedCode     int
		responseExpected string
		expectedError    []error
	}{
		{
			accountCreated:   &account.Account{ID: objectID, DocumentNumber: "12121212"},
			url:              account.AccountEndpoints.NewAccount,
			method:           "POST",
			bodyData:         "{\"documentNumber\": \"12121212\"}",
			expectedCode:     http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"documentNumber\":\"%v\"}}", http.StatusCreated, objectID.Hex(), "12121212"),
			expectedError:    nil,
		},
		{
			accountCreated:   nil,
			url:              account.AccountEndpoints.NewAccount,
			method:           "POST",
			bodyData:         "{\"documentNumber\": \"12121212\"}",
			expectedCode:     http.StatusBadRequest,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"%v\"]}", http.StatusBadRequest, "duplicate document"),
			expectedError:    []error{errors.New(helper.DuplicateMessageError)},
		},
		{
			accountCreated:   nil,
			url:              account.AccountEndpoints.NewAccount,
			method:           "POST",
			bodyData:         "{\"documentNumber\": \"\"}",
			expectedCode:     http.StatusBadRequest,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"%v\"]}", http.StatusBadRequest, "DocumentNumber field invalid"),
			expectedError:    nil,
		},
		{
			accountCreated:   nil,
			url:              account.AccountEndpoints.NewAccount,
			method:           "POST",
			bodyData:         "{\"documentNumber\": }",
			expectedCode:     http.StatusBadRequest,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"%v\"]}", http.StatusBadRequest, "DocumentNumber field invalid"),
			expectedError:    nil,
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockAccountService.EXPECT().NewAccount(gomock.Any()).Return(item.accountCreated, item.expectedError),
		)

		w := helper.ProcessRequest(suite.router, item.method, item.url, item.bodyData)

		suite.Assert().Equal(item.expectedCode, w.Code)
		suite.Assert().Equal(item.responseExpected, w.Body.String())
	}
}

func (suite *ControllerTestSuite) TestFindByID() {
	objectID := primitive.NewObjectID()
	data := []struct {
		expectedAccount  *account.Account
		url              string
		method           string
		bodyData         string
		expectedCode     int
		responseExpected string
		expectedError    []error
	}{
		{
			expectedAccount:  &account.Account{ID: objectID, DocumentNumber: "12121212"},
			url:              fmt.Sprintf("/accounts/%v", objectID.Hex()),
			method:           "GET",
			expectedCode:     http.StatusOK,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"documentNumber\":\"%v\"}}", http.StatusOK, objectID.Hex(), "12121212"),
			expectedError:    nil,
		},
		{
			expectedAccount:  nil,
			url:              fmt.Sprintf("/accounts/%v", objectID.Hex()),
			method:           "GET",
			expectedCode:     http.StatusNotFound,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"%v\"]}", http.StatusNotFound, "not found"),
			expectedError:    []error{errors.New(helper.NotFoundMessageError)},
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockAccountService.EXPECT().FindByID(gomock.Any()).Return(item.expectedAccount, item.expectedError),
		)
		suite.accountController = account.NewAccountController(suite.mockAccountService)

		w := helper.ProcessRequest(suite.router, item.method, item.url, item.bodyData)

		suite.Assert().Equal(item.expectedCode, w.Code)
		suite.Assert().Equal(item.responseExpected, w.Body.String())
	}
}
