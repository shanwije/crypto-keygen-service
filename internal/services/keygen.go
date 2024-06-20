package services

import (
	"crypto-keygen-service/internal/network_factory"
	"crypto-keygen-service/internal/network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/network_factory/generators/ethereum"
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/util/encryption"
	"crypto-keygen-service/internal/util/errors"
	log "github.com/sirupsen/logrus"
)

type KeyGenService struct {
	generators map[string]network_factory.KeyGenerator
	repository repositories.KeyGenRepository
}

func NewKeyGenService(repo repositories.KeyGenRepository, masterSeed []byte) *KeyGenService {
	service := &KeyGenService{
		generators: make(map[string]network_factory.KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{MasterSeed: masterSeed})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{MasterSeed: masterSeed})
	// Add more networks here
	return service
}

func (s *KeyGenService) RegisterGenerator(network string, generator network_factory.KeyGenerator) {
	s.generators[network] = generator
}

func (s *KeyGenService) GetKeysAndAddress(userID int, network string) (network_factory.KeyPairAndAddress, error) {
	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Request to get keys and address")

	exists, err := s.repository.KeyExists(userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to check if keys exist")
		return network_factory.KeyPairAndAddress{}, err
	}

	if exists {
		return s.retrieveExistingKeys(userID, network)
	}

	return s.generateAndSaveKeys(userID, network)
}

func (s *KeyGenService) retrieveExistingKeys(userID int, network string) (network_factory.KeyPairAndAddress, error) {
	keyData, err := s.repository.GetKey(userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve existing keys")
		return network_factory.KeyPairAndAddress{}, err
	}

	privateKey, err := encryption.Decrypt(keyData.EncryptedPrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt private key")
		return network_factory.KeyPairAndAddress{}, err
	}

	return network_factory.KeyPairAndAddress{
		Address:    keyData.Address,
		PublicKey:  keyData.PublicKey,
		PrivateKey: privateKey,
	}, nil
}

func (s *KeyGenService) generateAndSaveKeys(userID int, network string) (network_factory.KeyPairAndAddress, error) {
	generator, exists := s.generators[network]
	if !exists {
		log.WithFields(log.Fields{
			"network": network,
		}).Error("Unsupported network")
		return network_factory.KeyPairAndAddress{}, errors.ErrUnsupportedNetwork
	}

	keyPairAndAddress, err := generator.GenerateKeyPairAndAddress(userID)
	if err != nil {
		log.WithError(err).Error("Failed to generate key pair")
		return network_factory.KeyPairAndAddress{}, err
	}

	encryptedPrivateKey, err := encryption.Encrypt(keyPairAndAddress.PrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt private key")
		return network_factory.KeyPairAndAddress{}, err
	}

	keyData := repositories.KeyData{
		UserID:              userID,
		Network:             network,
		Address:             keyPairAndAddress.Address,
		PublicKey:           keyPairAndAddress.PublicKey,
		EncryptedPrivateKey: encryptedPrivateKey,
	}

	err = s.repository.SaveKey(userID, network, keyData)
	if err != nil {
		log.WithError(err).Error("Failed to save keys")
		return network_factory.KeyPairAndAddress{}, err
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Successfully generated and saved keys")

	return keyPairAndAddress, nil
}
