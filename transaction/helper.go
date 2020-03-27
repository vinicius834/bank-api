package transaction

import (
	"fmt"
	"log"
	"net/http"
)

const (
	TransactionCollection     = "Transaction"
	OperationTypeCollection   = "OperationType"
	InvalidAccountError       = "Account is invalid"
	InvalidOperationTypeError = "OperationType is invalid"
)

type Helper struct {
	httpClient *http.Client
}

type IHelper interface {
	AccountValidator(accountId string) bool
	TransformAmount(amount float64, isCredit bool) float64
}

func NewTransactionHelper(httpClient *http.Client) *Helper {
	return &Helper{httpClient}
}

func (helper *Helper) AccountValidator(accountId string) bool {
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/accounts/%v", accountId))
	if err != nil {
		log.Println(err)
		return false
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

func (helper *Helper) TransformAmount(amount float64, isCredit bool) float64 {
	if !isCredit {
		return amount * (-1)
	}
	return amount
}
