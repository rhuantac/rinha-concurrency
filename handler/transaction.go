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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRequest struct {
	Value           int                   `json:"valor" binding:"required,numeric"`
	TransactionType model.TransactionType `json:"tipo" binding:"required,oneof=c d"`
	Description     string                `json:"descricao" binding:"required,min=1,max=10"`
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
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "ID inválido."})
			return
		}

		var req TransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Message: "Payload inválido."})
			return
		}

		if req.TransactionType == model.DebitTransaction {
			req.Value *= -1
		}

		usersColl := db.Collection("users")
		filter := bson.D{{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$inc", Value: bson.D{{Key: "current_balance", Value: req.Value}}}}

		userCount, err := usersColl.CountDocuments(c, filter)
		if userCount < 1 || err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Cliente não encontrado."})
			return
		}

		updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var user model.User
		result := usersColl.FindOneAndUpdate(c, filter, update, updateOpts)
		result.Decode(&user)
		if result.Err() != nil {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Message: "Saldo insuficiente."})
			return
		}

		tx := model.Transaction{
			Value:           abs(req.Value),
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

		c.JSON(http.StatusOK, TransactionResponse{Limit: user.Limit, Balance: user.CurrentBalance})
	}
}

func abs(val int) int {
	if val < 0 {
		return -val
	}

	return val
}
