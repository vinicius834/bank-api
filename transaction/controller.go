package transaction

import (
	"bank-api/helper"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService ITransactionService
	helper             IHelper
}

type ITransactionController interface {
	NewTransaction(c *gin.Context)
	NewOperationType(c *gin.Context)
}

func NewTransactionController(transactionService ITransactionService, helper IHelper) *TransactionController {
	return &TransactionController{transactionService: transactionService, helper: helper}
}

func (transactionController *TransactionController) NewTransaction(c *gin.Context) {
	var newTransaction Transaction

	_ = json.NewDecoder(c.Request.Body).Decode(&newTransaction)

	errs := helper.ValidateFields(newTransaction)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	accountValid := transactionController.helper.AccountValidator(newTransaction.AccountID.Hex())
	if !accountValid {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, []error{errors.New(InvalidAccountError)}))
		return
	}

	transactionCreated, errs := transactionController.transactionService.NewTransaction(newTransaction)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	c.JSON(http.StatusCreated, helper.NewResponse(http.StatusCreated, transactionCreated))
}

func (transactionController *TransactionController) NewOperationType(c *gin.Context) {
	var newOperationType OperationType

	_ = json.NewDecoder(c.Request.Body).Decode(&newOperationType)

	errs := helper.ValidateFields(newOperationType)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	operationTypeCreated, errs := transactionController.transactionService.NewOperationType(newOperationType)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	c.JSON(http.StatusCreated, helper.NewResponse(http.StatusCreated, operationTypeCreated))
}
