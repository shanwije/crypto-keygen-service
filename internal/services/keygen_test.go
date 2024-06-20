package services

import (
	"context"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/ethereum"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndRetrieveKeyPair(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Clear the test DB
	client.Database("encryption-keygen-services-test").Drop(context.Background())

	repo := repositories.NewMongoRepository(client, "encryption-keygen-services-test", "currency_network_factory")
	keyService := NewKeyGenService(repo)
	keyService.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{})
	keyService.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{})

	userID := 1
	network := "bitcoin"

	// new key pair
	address, publicKey, privateKey, err := keyService.GetKeysAndAddress(userID, network)
	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.NotEmpty(t, publicKey)
	assert.NotEmpty(t, privateKey)

	// existing key pair
	address2, publicKey2, privateKey2, err := keyService.GetKeysAndAddress(userID, network)
	assert.NoError(t, err)
	assert.Equal(t, address, address2)
	assert.Equal(t, publicKey, publicKey2)
	assert.Equal(t, privateKey, privateKey2)
}
