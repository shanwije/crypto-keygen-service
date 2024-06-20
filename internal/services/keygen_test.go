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
	btcAddress1, btcPubKey1, btcPrivKey1, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")

	btcAddress2, btcPubKey2, btcPrivKey2, err := service.GetKeysAndAddress(userID, bitcoinNetwork)
	assert.NoError(t, err, "Expected no error for Bitcoin key generation")

	assert.Equal(t, btcAddress1, btcAddress2, "Expected same Bitcoin address for same user ID and network")
	assert.Equal(t, btcPubKey1, btcPubKey2, "Expected same Bitcoin public key for same user ID and network")
	assert.Equal(t, btcPrivKey1, btcPrivKey2, "Expected same Bitcoin private key for same user ID and network")

	// Test Ethereum consistency
	ethAddress1, ethPubKey1, ethPrivKey1, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")

	ethAddress2, ethPubKey2, ethPrivKey2, err := service.GetKeysAndAddress(userID, ethereumNetwork)
	assert.NoError(t, err, "Expected no error for Ethereum key generation")

	assert.Equal(t, ethAddress1, ethAddress2, "Expected same Ethereum address for same user ID and network")
	assert.Equal(t, ethPubKey1, ethPubKey2, "Expected same Ethereum public key for same user ID and network")
	assert.Equal(t, ethPrivKey1, ethPrivKey2, "Expected same Ethereum private key for same user ID and network")
}
