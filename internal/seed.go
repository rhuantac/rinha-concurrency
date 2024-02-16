package internal

import (
	"context"
	"os"

	"github.com/rhuantac/rinha-concurrency/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedDb(client *mongo.Client) {
	db := client.Database(os.Getenv("MONGO_DATABASE"))
	coll := db.Collection("users")
	users := []interface{}{
		model.User{ID: 1, Limit: 100000, InitialBalance: 0, CurrentBalance: 0},
		model.User{ID: 2, Limit: 80000, InitialBalance: 0, CurrentBalance: 0},
		model.User{ID: 3, Limit: 1000000, InitialBalance: 0, CurrentBalance: 0},
		model.User{ID: 4, Limit: 10000000, InitialBalance: 0, CurrentBalance: 0},
		model.User{ID: 5, Limit: 500000, InitialBalance: 0, CurrentBalance: 0},
	}
	_, err := coll.InsertMany(context.TODO(), users)

	if err != nil {
		panic(err)
	}
}

func ClearDb(client *mongo.Client) {
	client.Database(os.Getenv("MONGO_DATABASE")).Drop(context.TODO())
}
