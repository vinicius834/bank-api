package transaction_test

import (
	"bank-api/helper"
	"bank-api/transaction"
	mock_transaction "bank-api/transaction/mocks"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var accountId = "5e7aa59d62f361b30e620fea"

type ControllerTestSuite struct {
	suite.Suite
	transactionController     *transaction.TransactionController
	mockTransactionService    *mock_transaction.MockITransactionService
	mockTransactionRepository *mock_transaction.MockITransactionRepository
	mockTransactionHelper     *mock_transaction.MockIHelper
	mockCtrl                  *gomock.Controller
	router                    *gin.Engine
}

func (suite *ControllerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockTransactionService = mock_transaction.NewMockITransactionService(ctrl)
	suite.mockTransactionRepository = mock_transaction.NewMockITransactionRepository(ctrl)
	suite.mockTransactionHelper = mock_transaction.NewMockIHelper(ctrl)
	suite.initRouter()
	suite.transactionController = transaction.NewTransactionController(suite.mockTransactionService, suite.mockTransactionHelper)

	suite.setEndpoints()
}

func (suite *ControllerTestSuite) initRouter() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func (suite *ControllerTestSuite) setEndpoints() {
	v1 := suite.router.Group("")
	{
		v1.POST(transaction.TransactionEndpoints.NewTransaction, suite.transactionController.NewTransaction)
		v1.POST(transaction.TransactionEndpoints.NewOperationType, suite.transactionController.NewOperationType)
	}
}

func TestExampleSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (suite *ControllerTestSuite) TestNewTransaction() {
	objectID := primitive.NewObjectID()
	date := time.Now().Format(time.RFC3339Nano)
	data := []struct {
		transactionCreated *transaction.Transaction
		url                string
		method             string
		bodyData           string
		expectedCode       int
		responseExpected   string
		expectedError      []error
		accountValid       bool
	}{
		{
			transactionCreated: &transaction.Transaction{
				ID:            objectID,
				AccountID:     objectID,
				OperationType: objectID,
				Amount:        200.5,
				EventDate:     date,
			},
			url:          transaction.TransactionEndpoints.NewTransaction,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"accountID\": \"%v\", \"operationTypeID\": \"%v\", \"amount\": 200.5}", accountId, objectID.Hex()),
			expectedCode: http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"accountID\":\"%v\",\"operationTypeID\":\"%v\",\"amount\":200.5,\"eventDate\":\"%v\"}}",
				http.StatusCreated, objectID.Hex(), objectID.Hex(), objectID.Hex(), date),
			expectedError: nil,
			accountValid:  true,
		},
		{
			transactionCreated: &transaction.Transaction{
				ID:            objectID,
				AccountID:     objectID,
				OperationType: objectID,
				Amount:        -200.5,
				EventDate:     date,
			},
			url:          transaction.TransactionEndpoints.NewTransaction,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"accountID\": \"%v\", \"operationTypeID\": \"%v\", \"amount\": 200.5}", accountId, objectID.Hex()),
			expectedCode: http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"accountID\":\"%v\",\"operationTypeID\":\"%v\",\"amount\":-200.5,\"eventDate\":\"%v\"}}",
				http.StatusCreated, objectID.Hex(), objectID.Hex(), objectID.Hex(), date),
			expectedError: nil,
			accountValid:  true,
		},
		{
			transactionCreated: &transaction.Transaction{
				ID:            objectID,
				AccountID:     objectID,
				OperationType: objectID,
				Amount:        -200.5,
				EventDate:     date,
			},
			url:          transaction.TransactionEndpoints.NewTransaction,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"accountID\": \"%v\", \"operationTypeID\": \"%v\", \"amount\": 200.5}", "uwhswuhswhuswhhsuw", "uwhswuhswhuswhhsuw"),
			expectedCode: http.StatusBadRequest,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"AccountID field invalid\",\"OperationType field invalid\",\"Amount field invalid\"]}",
				http.StatusBadRequest),
			expectedError: nil,
			accountValid:  false,
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockTransactionHelper.EXPECT().AccountValidator(gomock.Any()).Return(item.accountValid),
			suite.mockTransactionService.EXPECT().NewTransaction(gomock.Any()).Return(item.transactionCreated, item.expectedError),
		)

		w := helper.ProcessRequest(suite.router, item.method, item.url, item.bodyData)

		suite.Assert().Equal(item.expectedCode, w.Code)
		suite.Assert().Equal(item.responseExpected, w.Body.String())
	}
}

func (suite *ControllerTestSuite) TestNewOperationType() {
	objectID := primitive.NewObjectID()
	descriptionCompraParcelada := "COMPRA PARCELADA"
	descriptionPagamento := "PAGAMENTO"

	data := []struct {
		operationTypeCreated *transaction.OperationType
		url                  string
		method               string
		bodyData             string
		expectedCode         int
		responseExpected     string
		expectedError        []error
	}{
		{
			operationTypeCreated: &transaction.OperationType{
				ID:          objectID,
				Description: descriptionCompraParcelada,
				IsCredit:    false,
			},
			url:          transaction.TransactionEndpoints.NewOperationType,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"description\": \"%v\", \"isCredit\": %v}", descriptionCompraParcelada, false),
			expectedCode: http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"description\":\"%v\",\"isCredit\":%v}}",
				http.StatusCreated, objectID.Hex(), descriptionCompraParcelada, false),
			expectedError: nil,
		},
		{
			operationTypeCreated: &transaction.OperationType{
				ID:          objectID,
				Description: descriptionPagamento,
				IsCredit:    true,
			},
			url:          transaction.TransactionEndpoints.NewOperationType,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"description\": \"%v\", \"isCredit\": %v}", descriptionPagamento, true),
			expectedCode: http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"description\":\"%v\",\"isCredit\":%v}}",
				http.StatusCreated, objectID.Hex(), descriptionPagamento, true),
			expectedError: nil,
		},
		{
			operationTypeCreated: &transaction.OperationType{
				ID:          objectID,
				Description: descriptionPagamento,
				IsCredit:    true,
			},
			url:          transaction.TransactionEndpoints.NewOperationType,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"description\": \"\", \"isCredit\": %v}", true),
			expectedCode: http.StatusBadRequest,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"errors\":[\"%v\"]}",
				http.StatusBadRequest, "Description field invalid"),
			expectedError: nil,
		},
		{
			operationTypeCreated: &transaction.OperationType{
				ID:          objectID,
				Description: descriptionPagamento,
				IsCredit:    false,
			},
			url:          transaction.TransactionEndpoints.NewOperationType,
			method:       "POST",
			bodyData:     fmt.Sprintf("{\"description\":\"%v\"}", descriptionPagamento),
			expectedCode: http.StatusCreated,
			responseExpected: fmt.Sprintf("{\"status\":%v,\"data\":{\"id\":\"%v\",\"description\":\"%v\",\"isCredit\":%v}}",
				http.StatusCreated, objectID.Hex(), descriptionPagamento, true),
			expectedError: nil,
		},
	}

	for _, item := range data {
		gomock.InOrder(
			suite.mockTransactionService.EXPECT().NewOperationType(gomock.Any()).Return(item.operationTypeCreated, item.expectedError),
		)

		w := helper.ProcessRequest(suite.router, item.method, item.url, item.bodyData)

		suite.Assert().Equal(item.expectedCode, w.Code)
		suite.Assert().Equal(item.responseExpected, w.Body.String())
	}
}
