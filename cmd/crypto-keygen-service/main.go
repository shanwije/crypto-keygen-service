package main

import (
	"context"
	"crypto-keygen-service/internal/db/mongo"
	"crypto-keygen-service/internal/util/encryption"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"crypto-keygen-service/internal/handlers"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	validateEnv()

	mongoURI := os.Getenv("MONGODB_URI")
	serverPort := os.Getenv("SERVER_PORT")
	dbName := os.Getenv("DB_NAME")
	dbCollection := os.Getenv("DB_COLLECTION")
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	masterSeed := os.Getenv("MASTER_SEED")

	setupEncryption(encryptionKey)
	database := setupDatabase(mongoURI, dbName, dbCollection)

	keyGenRepository := repositories.NewKeyGenRepository(database)
	keyGenService := services.NewKeyGenService(keyGenRepository, []byte(masterSeed))
	keyGenHandler := handlers.NewKeyGenHandler(keyGenService)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		healthCheck(c, database)
	})
	keyGenHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func validateEnv() {
	requiredVars := []string{"MONGODB_URI", "SERVER_PORT", "DB_NAME", "DB_COLLECTION", "ENCRYPTION_KEY", "MASTER_SEED"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}
}

func setupEncryption(key string) {
	if err := encryption.Setup(key); err != nil {
		log.Fatalf("Error setting up encryption: %v", err)
	}
}

func setupDatabase(uri, dbName, collection string) *mongo.MongoDatabase {
	database, err := mongo.NewMongoDatabase(uri, dbName, collection)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	return database
}

// todo: update health check for generic db
func healthCheck(c *gin.Context, database *mongo.MongoDatabase) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := database.Client.Ping(ctx, readpref.Primary()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
