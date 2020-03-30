package transaction

import (
	"bank-api/helper"
	"bank-api/storage"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionRepository struct {
	Database                storage.IMongoDB
	TransactionCollection   string
	OperationTypeCollection string
}

type ITransactionRepository interface {
	NewTransaction(newTransaction Transaction) (*Transaction, []error)
	FindOperationTypeById(id string) (*OperationType, []error)
	NewOperationType(newOperationType OperationType) (*OperationType, []error)
}

func NewTransactionRepository(mongoDB storage.IMongoDB, transactionCollection, operationTypeCollection string) *TransactionRepository {
	return &TransactionRepository{
		Database:                mongoDB,
		TransactionCollection:   transactionCollection,
		OperationTypeCollection: operationTypeCollection,
	}
}

func (transactionRepository *TransactionRepository) NewTransaction(newTransaction Transaction) (*Transaction, []error) {
	newTransaction.ID = primitive.NewObjectID()
	_, err := transactionRepository.Database.InsertOne(transactionRepository.TransactionCollection, newTransaction)

	if err != nil {
		return nil, []error{err}
	}
	return &newTransaction, nil
}

func (transactionRepository *TransactionRepository) NewOperationType(newOperationType OperationType) (*OperationType, []error) {
	newOperationType.ID = primitive.NewObjectID()
	_, err := transactionRepository.Database.InsertOne(transactionRepository.OperationTypeCollection, newOperationType)

	if err != nil {
		return nil, []error{err}
	}
	return &newOperationType, nil
}

func (transactionRepository *TransactionRepository) FindOperationTypeById(id string) (*OperationType, []error) {
	var operationTypeFound *OperationType
	idFormatted, _ := primitive.ObjectIDFromHex(id)
	fmt.Println(idFormatted)
	filter := bson.M{"_id": idFormatted}
	response := transactionRepository.Database.FindOne(transactionRepository.OperationTypeCollection, filter)
	err := response.Decode(&operationTypeFound)
	fmt.Println(err)
	if err != nil {
		return nil, []error{errors.New(helper.NotFoundMessageError)}
	}
	return operationTypeFound, nil
}
