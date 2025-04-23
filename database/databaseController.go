package database

import(
	"fmt"
	"log"
	"time"
	"os"
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	)

func DbInstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the env")
	}

	MongoDB := os.Getenv("MONGODB_URL")

// mongo.NewClient(options.Client().ApplyURI(MongoDB))
	// Create client with Connect (recommended)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	if err!= nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}
	fmt.Println("Database connected")
	log.Println("Connected to MongoDB!")
	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}