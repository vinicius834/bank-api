package storage_test

import (
	"bank-api/helper"
	"bank-api/storage"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDBTestSuite struct {
	suite.Suite
	mongoDB storage.IMongoDB
}

func (suite *MongoDBTestSuite) SetupTest() {
	suite.mongoDB = storage.NewMongoDB(helper.MongoDBUrlTest, helper.DBNameTest)
}

func (suite *MongoDBTestSuite) TearDownTest() {
	suite.mongoDB.GetDatabase(helper.DBNameTest).Drop(context.Background())
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoDBTestSuite))
}

func (suite *MongoDBTestSuite) TestNewAccount() {
	document := bson.M{"name": "Vinicius"}

	response, err := suite.mongoDB.InsertOne("TestCollection", document)
	if err != nil {
		suite.T().Errorf(fmt.Sprint(err))
	}
	suite.Assert().NotNil(response)
}

func (suite *MongoDBTestSuite) TestFindOne() {
	document := bson.M{"name": "Vinicius"}
	response, err := suite.mongoDB.InsertOne("TestCollection", document)

	filter := bson.M{"_id": response.InsertedID}
	documentFound := suite.mongoDB.FindOne("TestCollection", filter)

	if err = documentFound.Err(); err != nil {
		suite.T().Errorf(fmt.Sprint(err))
	}
	suite.Assert().NotNil(response)
}
