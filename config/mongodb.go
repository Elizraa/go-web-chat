package config

import (
	"context"
	"log"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient is exported Mongo Database client
var MongoDBClient *mongo.Client

// ConnectDatabase is used to connect the MongoDB database
func ConnectMongoDatabase() {
	log.Println("Mongo Database connecting...")
	uri := "mongodb+srv://" + url.QueryEscape(os.Getenv("MONGODB_USER")) + ":" +
		url.QueryEscape(os.Getenv("MONGODB_PASSWORD")) + "@" + os.Getenv("MONGODB_CLUSTER") +
		"/?retryWrites=true&w=majority"
	log.Println(uri)
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	MongoDBClient = client
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = MongoDBClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected.")
}
