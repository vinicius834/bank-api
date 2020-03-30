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
	CheckAccountHasLimit(id string, amountToCheck float64) []error
	UpdateLimit(id string, amount float64) []error
}

func NewAccountService(accountRepository IAccountRepository) *AccountService {
	return &AccountService{accountRepository: accountRepository}
}

func (accountService *AccountService) NewAccount(newAccount Account) (*Account, []error) {
	_, errs := accountService.accountRepository.FindByDocument(newAccount.Document)
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

func (accountService *AccountService) CheckAccountHasLimit(id string, amountToCheck float64) []error {
	accountFound, errs := accountService.accountRepository.FindByID(id)
	if helper.ErrorsExist(errs) {
		return errs
	}

	if (accountFound.AvalaibleLimit + amountToCheck) < 0 {
		return []error{errors.New(AccountHasNotLimitError)}
	}
	return nil
}

func (accountService *AccountService) UpdateLimit(id string, amount float64) []error {
	accountToUpdate, errs := accountService.accountRepository.FindByID(id)
	if errs != nil {
		return []error{errors.New(helper.NotFoundMessageError)}
	}

	accountToUpdate.AvalaibleLimit = accountToUpdate.AvalaibleLimit + amount

	_, errs = accountService.accountRepository.UpdateLimit(*accountToUpdate)
	if helper.ErrorsExist(errs) {
		return errs
	}
	return nil
}
