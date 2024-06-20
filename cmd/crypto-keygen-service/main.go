package main

import (
	"context"
	"crypto-keygen-service/internal/util/encryption"
	"log"
	"os"
	"time"

	"crypto-keygen-service/internal/handler"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
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
	masterSeed := os.Getenv("MASTER_SEED")
	if masterSeed == "" {
		log.Fatalf("MASTER_SEED environment variable is not set")
	}
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

	keyGenRepository := repositories.NewMongoRepository(client, dbName, dbCollection)
	keyGenService := services.NewKeyGenService(keyGenRepository, []byte(masterSeed))
	keyGenHandler := handler.NewKeyGenHandler(keyGenService)

	router := gin.Default()
	keyGenHandler.RegisterRoutes(router)
	router.Run(":" + serverPort)
}
