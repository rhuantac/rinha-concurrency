package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rhuantac/rinha-concurrency/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRequest struct {
	Value           int                   `json:"valor"`
	TransactionType model.TransactionType `json:"tipo"`
	Description     string                `json:"description"`
}
type TransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type ErrorResponse struct {
	Message string
}

func TransactionHandler(db *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		var req TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if req.TransactionType == model.DebitTransaction {
			req.Value *= -1
		}

		usersColl := db.Collection("users")
		filter := bson.D{{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$inc", Value: bson.D{{Key: "current_balance", Value: req.Value}}}}

		var user model.User
		result := usersColl.FindOne(c, filter)
		result.Decode(&user)

		if result.Err() != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Cliente não encontrado."})
			return
		}

		newBalance := user.CurrentBalance + req.Value

		//Balance cannot be lower than limit value
		if newBalance < user.Limit*-1 {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Message: "Saldo insuficiente."})
			return
		}

		_, err = usersColl.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Falha ao efetuar transação."})
			log.Printf("Error updating user %d", userId)
		}

		tx := model.Transaction{
			Value:           req.Value,
			TransactionType: req.TransactionType,
			Description:     req.Description,
			CreatedAt:       time.Now(),
			UserId:          userId,
		}

		txColl := db.Collection("transactions")
		_, err = txColl.InsertOne(c, tx)

		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Falha ao efetuar transação."})
			log.Printf("Error creating transaction %s", req.Description)
		}

		c.JSON(http.StatusOK, TransactionResponse{Limit: user.Limit, Balance: newBalance})
	}
}
