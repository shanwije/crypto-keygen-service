package main

import (
	"context"
	"crypto-keygen-service/internal/util/encryption"
	"log"
	"os"
	"time"

	"crypto-keygen-service/internal/controller"
	"crypto-keygen-service/internal/repository"
	"crypto-keygen-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	serverPort := os.Getenv("SERVER_PORT")
	dbName := os.Getenv("DB_NAME")
	dbCollection := os.Getenv("DB_COLLECTION")
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	err = encryption.Setup(encryptionKey)
	if err != nil {
		log.Fatalf("Error setting up encryption: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	repo := repository.NewMongoRepository(client, dbName, dbCollection)
	keyService := service.NewKeyService(repo)
	keyController := controller.NewKeyController(keyService)

	router := gin.Default()
	keyController.RegisterRoutes(router)
	router.Run(":" + serverPort)
}
