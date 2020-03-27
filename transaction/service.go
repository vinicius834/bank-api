package transaction

import (
	"bank-api/helper"
	"errors"
	"time"
)

type TransactionService struct {
	transactionRepository ITransactionRepository
	transactionHelper     IHelper
}

type ITransactionService interface {
	NewTransaction(newTransaction Transaction) (*Transaction, []error)
	NewOperationType(newOperationType OperationType) (*OperationType, []error)
}

func NewTransactionService(transactionRepository ITransactionRepository, transactionHelper IHelper) *TransactionService {
	return &TransactionService{transactionRepository: transactionRepository, transactionHelper: transactionHelper}
}

func (transactionService *TransactionService) NewTransaction(newTransaction Transaction) (*Transaction, []error) {
	operationType, errs := transactionService.transactionRepository.FindOperationTypeById(newTransaction.OperationType.Hex())
	if helper.ErrorsExist(errs) {
		return nil, []error{errors.New(InvalidOperationTypeError)}
	}

	newTransaction.EventDate = time.Now().Format(time.RFC3339Nano)
	newTransaction.Amount = transactionService.transactionHelper.TransformAmount(newTransaction.Amount, operationType.IsCredit)
	transactionInsertedResponse, errs := transactionService.transactionRepository.NewTransaction(newTransaction)
	if helper.ErrorsExist(errs) {
		return nil, errs
	}
	return transactionInsertedResponse, nil
}

func (transactionService *TransactionService) NewOperationType(newOperationType OperationType) (*OperationType, []error) {
	transactionInsertedResponse, errs := transactionService.transactionRepository.NewOperationType(newOperationType)
	if helper.ErrorsExist(errs) {
		return nil, errs
	}
	return transactionInsertedResponse, nil
}
