package account_test

import (
	"bank-api/account"
	"bank-api/config"
	"bank-api/helper"
	"bank-api/storage"
	"context"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type InitAccountModuleTestSuite struct {
	suite.Suite
	mongoDB storage.IMongoDB
	router  *gin.Engine
}

func (suite *InitAccountModuleTestSuite) SetupTest() {
	suite.mongoDB = storage.NewMongoDB(helper.MongoDBUrlTest, helper.DBNameTest)
	suite.router = config.InitRouter()
}

func (suite *InitAccountModuleTestSuite) TearDownTest() {
	suite.mongoDB.GetDatabase(helper.DBNameTest).Drop(context.Background())
}

func TestInitAccountModuleTestSuite(t *testing.T) {
	suite.Run(t, new(InitAccountModuleTestSuite))
}

func (suite *InitAccountModuleTestSuite) TestnitAccountModule() {
	endPoints := []string{account.AccountEndpoints.NewAccount, account.AccountEndpoints.FindByID}
	router := account.InitAccountModule(suite.mongoDB, suite.router)
	routesInfo := router.Routes()
	for i, routeInfo := range routesInfo {
		suite.Assert().Equal(endPoints[i], routeInfo.Path)
	}
}
