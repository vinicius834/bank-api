package account

import (
	"bank-api/storage"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	NewAccount string
	FindByID   string
}

var (
	accountRepository *AccountRepository
	accountService    *AccountService
	accountController *AccountController
	AccountEndpoints  = initAccountEndpoints()
)

func initAccountEndpoints() *Endpoints {
	return &Endpoints{
		NewAccount: "/accounts",
		FindByID:   "/accounts/:id",
	}
}

func InitAccountModule(mongoDB storage.IMongoDB, router *gin.Engine) *gin.Engine {
	initRepository(mongoDB)
	initService()
	initControllers()
	initAccountEndpoints()
	router = setEndpoints(router)
	return router
}

func initRepository(mongoDB storage.IMongoDB) {
	accountRepository = NewAccountRepository(mongoDB, AccountCollection)
}

func initService() {
	accountService = NewAccountService(accountRepository)
}

func initControllers() {
	accountController = NewAccountController(accountService)
}

func setEndpoints(router *gin.Engine) *gin.Engine {
	v1 := router.Group("")
	{
		v1.POST(AccountEndpoints.NewAccount, accountController.NewAccount)
		v1.GET(AccountEndpoints.FindByID, accountController.FindByID)
	}
	return router
}
