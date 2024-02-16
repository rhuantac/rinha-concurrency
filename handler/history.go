package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rhuantac/rinha-concurrency/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Balance struct {
	Total int       `json:"total"`
	Date  time.Time `json:"data_extrato"`
	Limit int       `json:"limite"`
}

type HistoryResponse struct {
	Balance Balance `json:"saldo"`
}

func HistoryHandler(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		usersColl := db.Collection("users")
		filter := bson.D{{Key: "_id", Value: userId}}

		var user model.User
		result := usersColl.FindOne(c, filter)
		result.Decode(&user)

		if result.Err() != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Cliente não encontrado."})
			return
		}

		c.JSON(http.StatusOK, HistoryResponse{Balance: Balance{Total: user.CurrentBalance, Date: time.Now(), Limit: user.Limit}})
	}
}
