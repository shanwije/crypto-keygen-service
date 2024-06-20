package services

import (
	"context"
	"crypto-keygen-service/internal/db"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/util/encryption"
	"crypto-keygen-service/internal/util/errors"
	. "crypto-keygen-service/internal/util/network_factory"
	"crypto-keygen-service/internal/util/network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/network_factory/generators/ethereum"
	log "github.com/sirupsen/logrus"
)

type KeyGenService struct {
	generators map[string]KeyGenerator
	repository *repositories.KeyGenRepository
}

func NewKeyGenService(repo *repositories.KeyGenRepository, masterSeed []byte) *KeyGenService {
	service := &KeyGenService{
		generators: make(map[string]KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{MasterSeed: masterSeed})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{MasterSeed: masterSeed})
	// Add more networks here
	return service
}

func (s *KeyGenService) RegisterGenerator(network string, generator KeyGenerator) {
	s.generators[network] = generator
}

// GetKeysAndAddress checks if a record exists for the given userID and network.
// If the record exists, it queries the database again to retrieve the record and return it.
// If the record does not exist, it creates a new record.
// This approach is used to avoid relying on error handling for control flow,
// providing clearer and more maintainable code.
func (s *KeyGenService) GetKeysAndAddress(userID int, network string) (KeyPairAndAddress, error) {
	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Request to get keys and address")

	ctx := context.Background()
	exists, err := s.repository.KeyExists(ctx, userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to check if keys exist")
		return KeyPairAndAddress{}, err
	}

	if exists {
		return s.retrieveExistingKeys(ctx, userID, network)
	}

	return s.generateAndSaveKeys(ctx, userID, network)
}

func (s *KeyGenService) retrieveExistingKeys(ctx context.Context, userID int, network string) (KeyPairAndAddress, error) {
	keyData, err := s.repository.GetKey(ctx, userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve existing keys")
		return KeyPairAndAddress{}, err
	}

	privateKey, err := encryption.Decrypt(keyData.EncryptedPrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt private key")
		return KeyPairAndAddress{}, err
	}

	return KeyPairAndAddress{
		Address:    keyData.Address,
		PublicKey:  keyData.PublicKey,
		PrivateKey: privateKey,
	}, nil
}

func (s *KeyGenService) generateAndSaveKeys(ctx context.Context, userID int, network string) (KeyPairAndAddress, error) {
	generator, exists := s.generators[network]
	if !exists {
		log.WithFields(log.Fields{
			"network": network,
		}).Error("Unsupported network")
		return KeyPairAndAddress{}, errors.ErrUnsupportedNetwork
	}

	keyPairAndAddress, err := generator.GenerateKeyPairAndAddress(userID)
	if err != nil {
		log.WithError(err).Error("Failed to generate key pair")
		return KeyPairAndAddress{}, err
	}

	encryptedPrivateKey, err := encryption.Encrypt(keyPairAndAddress.PrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt private key")
		return KeyPairAndAddress{}, err
	}

	keyData := db.KeyData{
		UserID:              userID,
		Network:             network,
		Address:             keyPairAndAddress.Address,
		PublicKey:           keyPairAndAddress.PublicKey,
		EncryptedPrivateKey: encryptedPrivateKey,
	}

	err = s.repository.SaveKey(ctx, keyData)
	if err != nil {
		log.WithError(err).Error("Failed to save keys")
		return KeyPairAndAddress{}, err
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Successfully generated and saved keys")

	return keyPairAndAddress, nil
}
