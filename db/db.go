package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set in environment")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "cinema_system"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	Client = client
	Database = client.Database(dbName)
	log.Println("✅ Connected to MongoDB successfully!")
}

func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}
