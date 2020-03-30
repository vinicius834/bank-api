package account

import (
	"bank-api/helper"
	"bank-api/storage"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountRepository struct {
	Database       storage.IMongoDB
	CollectionName string
}

type IAccountRepository interface {
	NewAccount(newAccount Account) (*Account, []error)
	FindByID(id string) (*Account, []error)
	FindByDocument(documentNumber string) (*Account, []error)
	UpdateLimit(accountToUpdate Account) (*Account, []error)
}

func NewAccountRepository(mongoDB storage.IMongoDB, collectionName string) *AccountRepository {
	return &AccountRepository{
		Database:       mongoDB,
		CollectionName: collectionName,
	}
}

func (accountRepository *AccountRepository) NewAccount(newAccount Account) (*Account, []error) {
	newAccount.ID = primitive.NewObjectID()
	_, err := accountRepository.Database.InsertOne(accountRepository.CollectionName, newAccount)

	if err != nil {
		return nil, []error{err}
	}

	return &newAccount, nil
}

func (accountRepository *AccountRepository) FindByID(id string) (*Account, []error) {
	var accountFound *Account
	idFormatted, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idFormatted}
	response := accountRepository.Database.FindOne(accountRepository.CollectionName, filter)
	err := response.Decode(&accountFound)

	if err != nil {
		return nil, []error{errors.New(helper.NotFoundMessageError)}
	}
	return accountFound, nil
}

func (accountRepository *AccountRepository) FindByDocument(document string) (*Account, []error) {
	var accountFound *Account
	filter := bson.M{"document": document}
	err := accountRepository.Database.FindOne(accountRepository.CollectionName, filter).Decode(&accountFound)
	if err != nil {
		return nil, []error{errors.New(helper.NotFoundMessageError)}
	}
	return accountFound, nil
}

func (accountRepository *AccountRepository) UpdateLimit(accountToUpdate Account) (*Account, []error) {
	id := bson.M{"_id": accountToUpdate.ID}

	limit := bson.M{
		"$set": bson.M{"available_limit": accountToUpdate.AvalaibleLimit},
	}
	_, err := accountRepository.Database.UpdateOne(accountRepository.CollectionName, id, limit)
	if err != nil {
		return nil, []error{err}
	}

	return &accountToUpdate, nil
}
