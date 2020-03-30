package transaction_test

import (
	"bank-api/account"
	mock_account "bank-api/account/mocks"
	"bank-api/helper"
	"bank-api/transaction"
	mock_transaction "bank-api/transaction/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceTestSuite struct {
	suite.Suite
	transactionService        transaction.ITransactionService
	mockAccountService        *mock_account.MockIAccountService
	mockTransactionRepository *mock_transaction.MockITransactionRepository
	mockCtrl                  *gomock.Controller
}

func (suite *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockTransactionRepository = mock_transaction.NewMockITransactionRepository(ctrl)
	suite.mockAccountService = mock_account.NewMockIAccountService(ctrl)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestNewTransaction() {
	objectID := primitive.NewObjectID()
	data := []struct {
		newTransaction      transaction.Transaction
		operationTypeFound  *transaction.OperationType
		expectedTransaction *transaction.Transaction
		expectedError       []error
		operationTypeError  []error
		checkLimitError     []error
		updateLimitError    []error
	}{
		{
			newTransaction: transaction.Transaction{
				AccountID:     primitive.NewObjectID(),
				OperationType: primitive.NewObjectID(),
				Amount:        300.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			operationTypeFound: &transaction.OperationType{
				ID:          primitive.NewObjectID(),
				Description: "COMPRA A VISTA",
				IsCredit:    false,
			},
			expectedTransaction: &transaction.Transaction{
				ID:            primitive.NewObjectID(),
				AccountID:     primitive.NewObjectID(),
				OperationType: primitive.NewObjectID(),
				Amount:        -300.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			expectedError:      nil,
			operationTypeError: nil,
			checkLimitError:    nil,
			updateLimitError:   nil,
		},
		{
			newTransaction: transaction.Transaction{
				AccountID:     primitive.NewObjectID(),
				OperationType: primitive.NewObjectID(),
				Amount:        300.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			operationTypeFound: &transaction.OperationType{
				ID:          primitive.NewObjectID(),
				Description: "PAGAMENTO",
				IsCredit:    true,
			},
			expectedTransaction: &transaction.Transaction{
				ID:            primitive.NewObjectID(),
				AccountID:     primitive.NewObjectID(),
				OperationType: primitive.NewObjectID(),
				Amount:        300.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			expectedError:      nil,
			operationTypeError: nil,
			checkLimitError:    nil,
			updateLimitError:   nil,
		},
		{
			newTransaction: transaction.Transaction{
				AccountID:     primitive.NewObjectID(),
				OperationType: objectID,
				Amount:        300.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			operationTypeFound: &transaction.OperationType{
				ID:          objectID,
				Description: "COMPRA A VISTA",
				IsCredit:    false,
			},
			expectedTransaction: nil,
			expectedError:       []error{errors.New(account.AccountHasNotLimitError)},
			operationTypeError:  nil,
			checkLimitError:     []error{errors.New(account.AccountHasNotLimitError)},
			updateLimitError:    nil,
		},
		{
			newTransaction: transaction.Transaction{
				AccountID:     primitive.NewObjectID(),
				OperationType: primitive.NewObjectID(),
				Amount:        100.0,
				EventDate:     time.Now().Format(time.RFC3339Nano),
			},
			operationTypeFound:  nil,
			expectedTransaction: nil,
			expectedError:       []error{errors.New(transaction.InvalidOperationTypeError)},
			operationTypeError:  []error{errors.New(helper.NotFoundMessageError)},
			checkLimitError:     nil,
			updateLimitError:    nil,
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockTransactionRepository.EXPECT().FindOperationTypeById(gomock.Any()).Return(item.operationTypeFound, item.operationTypeError),
			suite.mockAccountService.EXPECT().CheckAccountHasLimit(gomock.Any(), gomock.Any()).Return(item.checkLimitError),
			suite.mockAccountService.EXPECT().UpdateLimit(gomock.Any(), gomock.Any()).Return(item.updateLimitError),
			suite.mockTransactionRepository.EXPECT().NewTransaction(gomock.Any()).Return(item.expectedTransaction, item.expectedError),
		)
		suite.transactionService = transaction.NewTransactionService(suite.mockTransactionRepository, &transaction.Helper{}, suite.mockAccountService)
		transactionCreated, errs := suite.transactionService.NewTransaction(item.newTransaction)
		suite.Assert().Equal(item.expectedTransaction, transactionCreated)
		suite.Assert().Equal(item.expectedError, errs)
	}
}

func (suite *ServiceTestSuite) TestNewOperationType() {
	objectID := primitive.NewObjectID()
	objectID1 := primitive.NewObjectID()
	data := []struct {
		newOperationType      transaction.OperationType
		expectedOperationType *transaction.OperationType
		expectedError         []error
	}{
		{
			newOperationType: transaction.OperationType{
				ID:          objectID,
				Description: "PAGAMENTO",
				IsCredit:    true,
			},
			expectedOperationType: &transaction.OperationType{
				ID:          objectID,
				Description: "PAGAMENTO",
				IsCredit:    true,
			},
			expectedError: nil,
		},
		{
			newOperationType: transaction.OperationType{
				ID:          objectID1,
				Description: "COMPRA PARCELDA",
				IsCredit:    false,
			},
			expectedOperationType: &transaction.OperationType{
				ID:          objectID1,
				Description: "COMPRA PARCELDA",
				IsCredit:    false,
			},
			expectedError: nil,
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockTransactionRepository.EXPECT().NewOperationType(gomock.Any()).Return(item.expectedOperationType, item.expectedError),
		)
		suite.transactionService = transaction.NewTransactionService(suite.mockTransactionRepository, &transaction.Helper{}, suite.mockAccountService)
		operationTypeCreated, errs := suite.transactionService.NewOperationType(item.newOperationType)
		suite.Assert().Equal(item.expectedOperationType, operationTypeCreated)
		suite.Assert().Equal(item.expectedError, errs)
	}
}
