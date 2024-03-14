package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		println("something went wrong while loading env file in database")
		return
	}
	mongoURL = os.Getenv("MONGO_URL")
}

func CreateMongoClient() *mongo.Client {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		log.Fatal("something went wrong while creating mongo instance")
	}
	defer cancel()
	return client

}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("jwt").Collection(collectionName)
	return collection
}
