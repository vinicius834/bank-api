package account

import (
	"bank-api/helper"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService IAccountService
}

type IAccountController interface {
	NewAccount(c *gin.Context)
	FindByID(id string)
}

func NewAccountController(accountService IAccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

func (accountController *AccountController) NewAccount(c *gin.Context) {
	var newAccount Account

	_ = json.NewDecoder(c.Request.Body).Decode(&newAccount)

	errs := helper.ValidateFields(newAccount)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	accountCreated, errs := accountController.accountService.NewAccount(newAccount)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, errs))
		return
	}

	c.JSON(http.StatusCreated, helper.NewResponse(http.StatusCreated, accountCreated))
}

func (accountController *AccountController) FindByID(c *gin.Context) {
	id := c.Param("id")

	accountFound, errs := accountController.accountService.FindByID(id)
	if helper.ErrorsExist(errs) {
		c.JSON(http.StatusNotFound, helper.NewErrorResponse(http.StatusNotFound, errs))
		return
	}
	c.JSON(http.StatusOK, helper.NewResponse(http.StatusOK, accountFound))
}
