package services

import (
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/util/currency_network_factory"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/ethereum"
	"crypto-keygen-service/internal/util/encryption"
	"crypto-keygen-service/internal/util/errors"
	log "github.com/sirupsen/logrus"
)

type KeyGenService struct {
	generators map[string]currency_network_factory.KeyGenerator
	repository repositories.KeyGenRepository
}

func NewKeyGenService(repo repositories.KeyGenRepository, masterSeed []byte) *KeyGenService {
	service := &KeyGenService{
		generators: make(map[string]currency_network_factory.KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{MasterSeed: masterSeed})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{MasterSeed: masterSeed})
	// Add more networks here
	return service
}

func (s *KeyGenService) RegisterGenerator(network string, generator currency_network_factory.KeyGenerator) {
	s.generators[network] = generator
}

func (s *KeyGenService) GetKeysAndAddress(userID int, network string) (string, string, string, error) {
	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Request to get keys and address")

	exists, err := s.repository.KeyExists(userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to check if keys exist")
		return "", "", "", err
	}

	if exists {
		return s.retrieveExistingKeys(userID, network)
	}

	return s.generateAndSaveKeys(userID, network)
}

func (s *KeyGenService) retrieveExistingKeys(userID int, network string) (string, string, string, error) {
	address, publicKey, encryptedPrivateKey, err := s.repository.GetKey(userID, network)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve existing keys")
		return "", "", "", err
	}

	privateKey, err := encryption.Decrypt(encryptedPrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt private key")
		return "", "", "", err
	}

	return address, publicKey, privateKey, nil
}

func (s *KeyGenService) generateAndSaveKeys(userID int, network string) (string, string, string, error) {
	generator, exists := s.generators[network]
	if !exists {
		log.WithFields(log.Fields{
			"network": network,
		}).Error("Unsupported network")
		return "", "", "", errors.ErrUnsupportedNetwork
	}

	address, publicKey, privateKey, err := generator.GenerateKeyPair(userID)
	if err != nil {
		log.WithError(err).Error("Failed to generate key pair")
		return "", "", "", err
	}

	encryptedPrivateKey, err := encryption.Encrypt(privateKey)
	if err != nil {
		log.WithError(err).Error("Failed to encrypt private key")
		return "", "", "", err
	}

	err = s.repository.SaveKey(userID, network, address, publicKey, encryptedPrivateKey)
	if err != nil {
		log.WithError(err).Error("Failed to save keys")
		return "", "", "", err
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Successfully generated and saved keys")

	return address, publicKey, privateKey, nil
}
