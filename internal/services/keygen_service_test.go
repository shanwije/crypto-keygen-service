package services_test

import (
	"context"
	"crypto-keygen-service/internal/db"
	dbi "crypto-keygen-service/internal/db"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/services"
	"crypto-keygen-service/internal/util/encryption"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sampleMasterSeed = "sample-master-seed"

type InMemoryDatabase struct {
	data map[int]map[string]db.KeyData
}

func (db *InMemoryDatabase) CreateIndexes(ctx context.Context) error {
	return nil
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		data: make(map[int]map[string]db.KeyData),
	}
}

func (db *InMemoryDatabase) SaveKey(ctx context.Context, keyData db.KeyData) error {
	if _, ok := db.data[keyData.UserID]; !ok {
		db.data[keyData.UserID] = make(map[string]dbi.KeyData)
	}
	db.data[keyData.UserID][keyData.Network] = keyData
	return nil
}

func (db *InMemoryDatabase) GetKey(ctx context.Context, userID int, network string) (db.KeyData, error) {
	if userKeys, ok := db.data[userID]; ok {
		if keyData, ok := userKeys[network]; ok {
			return keyData, nil
		}
	}
	return dbi.KeyData{}, errors.New("key not found")
}

func (db *InMemoryDatabase) KeyExists(ctx context.Context, userID int, network string) (bool, error) {
	if userKeys, ok := db.data[userID]; ok {
		if _, ok := userKeys[network]; ok {
			return true, nil
		}
	}
	return false, nil
}

func TestGenerateAndRetrieveKeys(t *testing.T) {
	encryptionKey := "4GRrhM8ClnrSmCrDvyFzPKdkJF9NcRkKwxlmIrsYhx0="
	err := encryption.Setup(encryptionKey)
	assert.NoError(t, err)

	inMemoryDB := NewInMemoryDatabase()
	repo := repositories.NewKeyGenRepository(inMemoryDB)
	service := services.NewKeyGenService(repo, []byte(sampleMasterSeed))

	userID := 12345
	network := "bitcoin"

	// Generate keys
	result1, err := service.GetKeysAndAddress(userID, network)
	assert.NoError(t, err)
	assert.NotEmpty(t, result1.Address)
	assert.NotEmpty(t, result1.PublicKey)
	assert.NotEmpty(t, result1.PrivateKey)

	// Retrieve keys
	result2, err := service.GetKeysAndAddress(userID, network)
	assert.NoError(t, err)
	assert.Equal(t, result1, result2)
}
