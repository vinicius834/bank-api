package account_test

import (
	"bank-api/account"
	mock_account "bank-api/account/mocks"
	"bank-api/helper"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceTestSuite struct {
	suite.Suite
	accountService        account.IAccountService
	mockAccountRepository *mock_account.MockIAccountRepository
	mockCtrl              *gomock.Controller
}

func (suite *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockAccountRepository = mock_account.NewMockIAccountRepository(ctrl)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestNewAccount() {
	data := []struct {
		newAccount      account.Account
		accountFound    *account.Account
		expectedAccount *account.Account
		expectedError   []error
		errorToReturn   []error
	}{
		{
			newAccount:      account.Account{DocumentNumber: "121221"},
			accountFound:    nil,
			expectedAccount: &account.Account{ID: primitive.NewObjectID(), DocumentNumber: "121221"},
			expectedError:   nil,
			errorToReturn:   []error{errors.New(helper.NotFoundMessageError)}},
		{
			newAccount:      account.Account{DocumentNumber: "121221"},
			accountFound:    &account.Account{ID: primitive.NewObjectID(), DocumentNumber: "121221"},
			expectedAccount: nil,
			expectedError:   []error{errors.New(helper.DuplicateMessageError)},
			errorToReturn:   nil},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockAccountRepository.EXPECT().FindByDocument(gomock.Any()).Return(item.accountFound, item.errorToReturn),
			suite.mockAccountRepository.EXPECT().NewAccount(item.newAccount).Return(item.expectedAccount, item.expectedError),
		)
		suite.accountService = account.NewAccountService(suite.mockAccountRepository)
		accountCreated, errs := suite.accountService.NewAccount(item.newAccount)
		suite.Assert().Equal(item.expectedAccount, accountCreated)
		suite.Assert().Equal(item.expectedError, errs)
	}
}

func (suite *ServiceTestSuite) TestFindbyID() {
	objectID := primitive.NewObjectID()
	data := []struct {
		expectedAccount *account.Account
		expectedError   []error
	}{
		{
			expectedAccount: &account.Account{ID: objectID, DocumentNumber: "121221"},
			expectedError:   nil},
		{
			expectedAccount: nil,
			expectedError:   []error{errors.New(helper.NotFoundMessageError)}},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockAccountRepository.EXPECT().FindByID(gomock.Any()).Return(item.expectedAccount, item.expectedError),
		)
		suite.accountService = account.NewAccountService(suite.mockAccountRepository)
		accountFound, errs := suite.accountService.FindByID(objectID.Hex())
		suite.Assert().Equal(item.expectedAccount, accountFound)
		suite.Assert().Equal(item.expectedError, errs)
	}
}
