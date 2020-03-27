package transaction

import (
	"bank-api/storage"
	"net/http"

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
	initRepository(mongoDB)
	initService()
	initControllers()
	initTransactionEndpoints()
	router = setEndpointsToRouter(router)
	return router
}

func initHelper() {
	transactionHelper = NewTransactionHelper(&http.Client{})
}

func initRepository(mongoDB storage.IMongoDB) {
	transactionRepository = NewTransactionRepository(mongoDB, TransactionCollection, OperationTypeCollection)
}

func initService() {
	transactionService = NewTransactionService(transactionRepository, transactionHelper)
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
