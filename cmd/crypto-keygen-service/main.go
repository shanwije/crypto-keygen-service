package main

import (
	"context"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
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

	repo := repository.NewMongoRepository(client, "crypto-keygen-service", "currency_factory")

	keyService := service.NewKeyService(repo)

	router := gin.Default()
	controller.RegisterRoutes(router, keyService)
	router.Run(":8080")
}
