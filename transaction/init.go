package transaction

import (
	"bank-api/account"
	"bank-api/storage"

	"github.com/gin-gonic/gin"
)

type Endpoints struct {
	NewTransaction   string
	NewOperationType string
}

var (
	transactionRepository *TransactionRepository
	transactionService    *TransactionService
	transactionController *TransactionController
	accountService        *account.AccountService
	transactionHelper     *Helper
	TransactionEndpoints  = initTransactionEndpoints()
)

func initTransactionEndpoints() *Endpoints {
	return &Endpoints{
		NewTransaction:   "/transactions",
		NewOperationType: "/transactions/operation-type",
	}
}

func InitTrasnactionModule(mongoDB storage.IMongoDB, router *gin.Engine) *gin.Engine {
	initAccountService(mongoDB)
	initRepository(mongoDB)
	initService()
	initControllers()
	initTransactionEndpoints()
	router = setEndpointsToRouter(router)
	return router
}

func initAccountService(mongoDB storage.IMongoDB) {
	accountService = account.NewAccountService(account.NewAccountRepository(mongoDB, account.AccountCollection))
}
func initHelper() {
	transactionHelper = NewTransactionHelper()
}

func initRepository(mongoDB storage.IMongoDB) {
	transactionRepository = NewTransactionRepository(mongoDB, TransactionCollection, OperationTypeCollection)
}

func initService() {
	transactionService = NewTransactionService(transactionRepository, transactionHelper, accountService)
}

func initControllers() {
	transactionController = NewTransactionController(transactionService, transactionHelper)
}

func setEndpointsToRouter(router *gin.Engine) *gin.Engine {
	v1 := router.Group("")
	{
		v1.POST(TransactionEndpoints.NewTransaction, transactionController.NewTransaction)
		v1.POST(TransactionEndpoints.NewOperationType, transactionController.NewOperationType)
	}
	return router
}
