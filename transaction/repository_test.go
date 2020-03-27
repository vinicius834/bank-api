package transaction_test

import (
	"bank-api/helper"
	"bank-api/storage"
	"bank-api/transaction"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryTestSuite struct {
	suite.Suite
	transactionRepository transaction.ITransactionRepository
	mongoDB               storage.IMongoDB
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.mongoDB = storage.NewMongoDB(helper.MongoDBUrlTest, helper.DBNameTest)
	suite.transactionRepository = transaction.NewTransactionRepository(suite.mongoDB, transaction.TransactionCollection, transaction.OperationTypeCollection)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.mongoDB.GetDatabase(helper.DBNameTest).Drop(context.Background())
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestNewTransaction() {
	newTransaction := transaction.Transaction{
		OperationType: primitive.NewObjectID(),
		AccountID:     primitive.NewObjectID(),
		Amount:        300.0,
		EventDate:     time.Now().Format(time.RFC3339Nano),
	}
	transactionCreated, errs := suite.transactionRepository.NewTransaction(newTransaction)
	if helper.ErrorsExist(errs) {
		suite.T().Errorf(fmt.Sprint(errs))
	}
	suite.Assert().Equal(newTransaction.AccountID, transactionCreated.AccountID)
	suite.Assert().Equal(newTransaction.OperationType, transactionCreated.OperationType)
	suite.Assert().Equal(newTransaction.AccountID, transactionCreated.AccountID)
	suite.Assert().Equal(newTransaction.Amount, transactionCreated.Amount)
	suite.Assert().Equal(newTransaction.EventDate, transactionCreated.EventDate)
}

func (suite *RepositoryTestSuite) TestNewOperationType() {
	newOperationtype := transaction.OperationType{
		Description: "Pagamento",
		IsCredit:    true,
	}
	operationtypeCreated, errs := suite.transactionRepository.NewOperationType(newOperationtype)
	if helper.ErrorsExist(errs) {
		suite.T().Errorf(fmt.Sprint(errs))
	}
	suite.Assert().Equal(newOperationtype.Description, operationtypeCreated.Description)
	suite.Assert().Equal(newOperationtype.IsCredit, operationtypeCreated.IsCredit)
}

func (suite *RepositoryTestSuite) TestFindOperationTypeById() {
	objectID := primitive.NewObjectID()
	data := []struct {
		operationType transaction.OperationType
		expected      *transaction.OperationType
		errorMessage  interface{}
	}{
		{
			operationType: transaction.OperationType{ID: objectID, Description: "COMPRA A VISTA", IsCredit: false},
			expected:      nil, errorMessage: nil,
		},
		{
			operationType: transaction.OperationType{ID: primitive.NewObjectID(), Description: "COMPRA PARCELADA", IsCredit: false},
			expected:      nil, errorMessage: helper.NotFoundMessageError,
		},
	}
	for i, item := range data {
		_, err := suite.mongoDB.InsertOne(transaction.OperationTypeCollection, item.operationType)
		if err != nil {
			suite.T().Errorf("%v", err)
		}
		if i != len(data)-1 {
			item.operationType.ID = objectID
			item.expected = &item.operationType
		} else {
			item.operationType.ID = primitive.NewObjectID()
		}
		operationTypeFound, errs := suite.transactionRepository.FindOperationTypeById(item.operationType.ID.Hex())
		if helper.ErrorsExist(errs) {
			suite.Assert().Equal(item.errorMessage, errs[0].Error())
		}
		suite.Assert().Equal(item.expected, operationTypeFound)
		suite.Assert().Equal(item.expected, operationTypeFound)
	}
}
