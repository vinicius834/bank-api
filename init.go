package main

import (
	"bank-api/account"
	"bank-api/config"
	"bank-api/helper"
	"bank-api/storage"
	"bank-api/transaction"

	"github.com/gin-gonic/gin"
)

var (
	router  *gin.Engine
	mongoDB storage.IMongoDB
)

type Person struct {
	name string
}

func main() {
	mongoDB = initDatabase()
	router = initRouter()
	initModules()
	run()
}

func initDatabase() storage.IMongoDB {
	return storage.NewMongoDB(helper.MongoDBUrl, helper.DBName)
}

func initRouter() *gin.Engine {
	return config.InitRouter()
}

func initModules() {
	router = account.InitAccountModule(mongoDB, router)
	router = transaction.InitTrasnactionModule(mongoDB, router)
}

func run() {
	router.Run()
}
