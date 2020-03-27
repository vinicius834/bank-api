package account_test

import (
	"bank-api/account"
	"bank-api/helper"
	"bank-api/storage"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryTestSuite struct {
	suite.Suite
	accountRepository account.IAccountRepository
	mongoDB           storage.IMongoDB
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.mongoDB = storage.NewMongoDB(helper.MongoDBUrlTest, helper.DBNameTest)
	suite.accountRepository = account.NewAccountRepository(suite.mongoDB, account.AccountCollection)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.mongoDB.GetDatabase(helper.DBNameTest).Drop(context.Background())
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestNewAccount() {
	newAccount := account.Account{DocumentNumber: "1234"}
	accountCreated, err := suite.accountRepository.NewAccount(newAccount)
	if err != nil {
		suite.T().Errorf(fmt.Sprint(err))
	}
	suite.Assert().Equal(newAccount.DocumentNumber, accountCreated.DocumentNumber)
}

func (suite *RepositoryTestSuite) TestFindByID() {
	data := []struct {
		acc          account.Account
		expected     *account.Account
		errorMessage interface{}
	}{
		{acc: account.Account{DocumentNumber: "12121"}, expected: nil, errorMessage: nil},
		{acc: account.Account{DocumentNumber: "12121221212121212"}, expected: nil, errorMessage: helper.NotFoundMessageError},
	}
	for i, item := range data {
		response, err := suite.accountRepository.NewAccount(item.acc)
		if err != nil {
			suite.T().Errorf("%v", err)
		}
		if i != len(data)-1 {
			item.acc.ID = response.ID
			item.expected = &item.acc
		} else {
			item.acc.ID = primitive.NewObjectID()
		}
		accountFound, errs := suite.accountRepository.FindByID(item.acc.ID.Hex())
		if helper.ErrorsExist(errs) {
			suite.Assert().Equal(item.errorMessage, errs[0].Error())
		}
		suite.Assert().Equal(item.expected, accountFound)
	}
}

func (suite *RepositoryTestSuite) TestFindByDocument() {
	data := []struct {
		acc          account.Account
		expected     *account.Account
		errorMessage interface{}
	}{
		{acc: account.Account{DocumentNumber: "12121"}, expected: nil, errorMessage: nil},
		{acc: account.Account{DocumentNumber: "12121221212121212"}, expected: nil, errorMessage: helper.NotFoundMessageError},
	}
	for i, item := range data {
		response, err := suite.accountRepository.NewAccount(item.acc)
		if err != nil {
			suite.T().Errorf("%v", err)
		}
		if i != len(data)-1 {
			item.acc.ID = response.ID
			item.expected = &item.acc
		} else {
			item.acc.DocumentNumber = "-91829198291982s"
		}
		accountFound, errs := suite.accountRepository.FindByDocument(item.acc.DocumentNumber)
		if helper.ErrorsExist(errs) {
			suite.Assert().Equal(item.errorMessage, errs[0].Error())
		}
		suite.Assert().Equal(item.expected, accountFound)
	}
}
