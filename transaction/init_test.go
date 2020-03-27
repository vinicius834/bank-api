package transaction_test

import (
	"bank-api/config"
	"bank-api/helper"
	"bank-api/storage"
	"bank-api/transaction"
	"context"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type InitTransactionModuleTestSuite struct {
	suite.Suite
	mongoDB storage.IMongoDB
	router  *gin.Engine
}

func (suite *InitTransactionModuleTestSuite) SetupTest() {
	suite.mongoDB = storage.NewMongoDB(helper.MongoDBUrlTest, helper.DBNameTest)
	suite.router = config.InitRouter()
}

func (suite *InitTransactionModuleTestSuite) TearDownTest() {
	suite.mongoDB.GetDatabase(helper.DBNameTest).Drop(context.Background())
}

func TestInitTransactionModuleTestSuite(t *testing.T) {
	suite.Run(t, new(InitTransactionModuleTestSuite))
}

func (suite *InitTransactionModuleTestSuite) TestnitTransactionModule() {
	endPoints := []string{
		transaction.TransactionEndpoints.NewTransaction,
		transaction.TransactionEndpoints.NewOperationType}
	router := transaction.InitTrasnactionModule(suite.mongoDB, suite.router)
	routesInfo := router.Routes()
	for i, routeInfo := range routesInfo {
		suite.Assert().Equal(endPoints[i], routeInfo.Path)
	}
}
