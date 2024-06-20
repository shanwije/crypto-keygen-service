package main

import (
	"crypto-keygen-service/internal/db/mongo"
	"crypto-keygen-service/internal/util/encryption"
	"log"
	"os"

	"crypto-keygen-service/internal/handlers"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	mongoDatabase, err := mongo.NewMongoDatabase(mongoURI, dbName, dbCollection)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	keyGenRepository := repositories.NewKeyGenRepository(mongoDatabase)
	keyGenService := services.NewKeyGenService(keyGenRepository, []byte(masterSeed))
	keyGenHandler := handlers.NewKeyGenHandler(keyGenService)

	router := gin.Default()
	keyGenHandler.RegisterRoutes(router)
	router.Run(":" + serverPort)
}
