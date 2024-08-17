package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
}

func (database *Database) New(client *mongo.Client) {
	database.Client = client
}

var database Database

func InitDb() {
	client, err := connectToMongodb()
	if err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
	database.New(client)
}

func connectToMongodb() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(DB_URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetDatabase() *mongo.Database {
	return database.Client.Database(DB_NAME)
}
