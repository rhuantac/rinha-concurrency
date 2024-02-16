package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongo() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_HOST"))
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return client
}

func DisconnectMongo(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
