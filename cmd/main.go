package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rhuantac/rinha-concurrency/config"
	"github.com/rhuantac/rinha-concurrency/handler"
	"github.com/rhuantac/rinha-concurrency/internal"
)

func init() {
	config.SetupEnvs()
}

func main() {	
	mongoClient := config.SetupMongo()
	defer config.DisconnectMongo(mongoClient)
	internal.ClearDb(mongoClient)
	internal.SeedDb(mongoClient)
	router := gin.Default()
	api := router.Group("/clientes/:id")

	api.POST("/transacoes", handler.TransactionHandler(nil))
	router.Run(":3000")
	
}
