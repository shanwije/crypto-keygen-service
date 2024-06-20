package services_test

import (
	"context"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
	"crypto-keygen-service/internal/util/encryption"
	"github.com/stretchr/testify/assert"
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

func setupMongoRepository() *repositories.MongoRepository {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return repositories.NewMongoRepository(client, "crypto-keygen-service", "crypto-wallet-service")
}

func TestServiceConsistency(t *testing.T) {
	repo := setupMongoRepository()
	service := services.NewKeyGenService(repo, []byte(sampleMasterSeed))

	userID := 12345
	bitcoinNetwork := "bitcoin"
	ethereumNetwork := "ethereum"

	// Test Bitcoin consistency
	btcKeys1, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")

	btcKeys2, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")

	assert.Equal(t, btcKeys1.Address, btcKeys2.Address, "Expected same Bitcoin address for same user ID and network")
	assert.Equal(t, btcKeys1.PublicKey, btcKeys2.PublicKey, "Expected same Bitcoin public key for same user ID and network")
	assert.Equal(t, btcKeys1.PrivateKey, btcKeys2.PrivateKey, "Expected same Bitcoin private key for same user ID and network")

	// Test Ethereum consistency
	ethKeys1, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")

	ethKeys2, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")

	assert.Equal(t, ethKeys1.Address, ethKeys2.Address, "Expected same Ethereum address for same user ID and network")
	assert.Equal(t, ethKeys1.PublicKey, ethKeys2.PublicKey, "Expected same Ethereum public key for same user ID and network")
	assert.Equal(t, ethKeys1.PrivateKey, ethKeys2.PrivateKey, "Expected same Ethereum private key for same user ID and network")
}
