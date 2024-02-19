package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rhuantac/rinha-concurrency/config"
	"github.com/rhuantac/rinha-concurrency/handler"
	"github.com/rhuantac/rinha-concurrency/internal"
)

func init() {
	config.SetupEnvs()
}

func main() {
	log.Print("Initializing server")
	mongoClient := config.SetupMongo()
	defer config.DisconnectMongo(mongoClient)
	internal.ClearDb(mongoClient)
	internal.SeedDb(mongoClient)

	redisClient := config.SetupRedis()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/clientes/:id")

	api.POST("/transacoes", handler.TransactionHandler(mongoClient.Database(os.Getenv("MONGO_DATABASE")), redisClient))
	api.GET("/extrato", handler.HistoryHandler(mongoClient.Database(os.Getenv("MONGO_DATABASE"))))
	log.Print("Running")
	router.Run()

}
