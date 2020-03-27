package account

import (
	"bank-api/helper"
	"errors"
)

type AccountService struct {
	accountRepository IAccountRepository
}

type IAccountService interface {
	NewAccount(newAccount Account) (*Account, []error)
	FindByID(id string) (*Account, []error)
}

func NewAccountService(accountRepository IAccountRepository) *AccountService {
	return &AccountService{accountRepository: accountRepository}
}

func (accountService *AccountService) NewAccount(newAccount Account) (*Account, []error) {
	_, errs := accountService.accountRepository.FindByDocument(newAccount.DocumentNumber)
	if errs == nil {
		return nil, []error{errors.New(helper.DuplicateMessageError)}
	}

	accountInsertedResponse, errs := accountService.accountRepository.NewAccount(newAccount)
	if helper.ErrorsExist(errs) {
		return nil, errs
	}
	return accountInsertedResponse, nil
}

func (accountService *AccountService) FindByID(id string) (*Account, []error) {
	accountFound, errs := accountService.accountRepository.FindByID(id)
	if helper.ErrorsExist(errs) {
		return nil, errs
	}
	return accountFound, nil
}
