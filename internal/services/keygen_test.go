package services_test

import (
	"context"
	"crypto-keygen-service/internal/db"
	mongoDB "crypto-keygen-service/internal/db/mongo"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
	"crypto-keygen-service/internal/util/encryption"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

const sampleEncryptionKey = "4GRrhM8ClnrSmCrDvyFzPKdkJF9NcRkKwxlmIrsYhx0="
const sampleMasterSeed = "sample-master-seed"

func TestMain(m *testing.M) {
	err := encryption.Setup(sampleEncryptionKey)
	if err != nil {
		panic("Failed to set up encryption: " + err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

func setupDatabase() db.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	mongoDatabase, err := mongoDB.NewMongoDatabase("mongodb://localhost:27017", "crypto-keygen-service-test", "crypto-wallet-service")
	if err != nil {
		panic(err)
	}

	return mongoDatabase
}

func setupKeyGenService(database db.Database) *services.KeyGenService {
	repo := repositories.NewKeyGenRepository(database)
	return services.NewKeyGenService(repo, []byte(sampleMasterSeed))
}

func TestServiceIntegration(t *testing.T) {
	database := setupDatabase()
	service := setupKeyGenService(database)

	userID := 12345
	bitcoinNetwork := "bitcoin"
	ethereumNetwork := "ethereum"

	// Clean up
	_, _ = database.(*mongoDB.MongoDatabase).Collection.DeleteMany(context.Background(), bson.M{"user_id": userID})

	// Test Bitcoin key generation and retrieval
	btcResult1, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")
	assert.NotEmpty(t, btcResult1.Address, "Expected non-empty Bitcoin address")
	assert.NotEmpty(t, btcResult1.PublicKey, "Expected non-empty Bitcoin public key")
	assert.NotEmpty(t, btcResult1.PrivateKey, "Expected non-empty Bitcoin private key")

	btcResult2, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")
	assert.Equal(t, btcResult1, btcResult2, "Expected same keys for repeated Bitcoin key generation")

	// Test Ethereum key generation and retrieval
	ethResult1, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")
	assert.NotEmpty(t, ethResult1.Address, "Expected non-empty Ethereum address")
	assert.NotEmpty(t, ethResult1.PublicKey, "Expected non-empty Ethereum public key")
	assert.NotEmpty(t, ethResult1.PrivateKey, "Expected non-empty Ethereum private key")

	ethResult2, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")
	assert.Equal(t, ethResult1, ethResult2, "Expected same keys for repeated Ethereum key generation")

	_, _ = database.(*mongoDB.MongoDatabase).Collection.DeleteMany(context.Background(), bson.M{"user_id": userID})
}
