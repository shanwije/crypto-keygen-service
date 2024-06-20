package service

import (
	"log"

	"crypto-keygen-service/internal/repository"
	"crypto-keygen-service/internal/util/currency_network_factory"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/ethereum"
	"crypto-keygen-service/internal/util/encryption"
	"crypto-keygen-service/internal/util/errors"
)

type KeyService struct {
	generators map[string]currency_network_factory.KeyGenerator
	repository repository.Repository
}

func NewKeyService(repo repository.Repository) *KeyService {
	service := &KeyService{
		generators: make(map[string]currency_network_factory.KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{})
	// Add more networks here
	return service
}

func (s *KeyService) RegisterGenerator(network string, generator currency_network_factory.KeyGenerator) {
	s.generators[network] = generator
}

// GetKeysAndAddress retrieves or generates keys and address for a user on a specific network.
func (s *KeyService) GetKeysAndAddress(userID int, network string) (string, string, string, error) {
	// Check if keys already exist in the repository
	exists, err := s.repository.KeyExists(userID, network)
	if err != nil {
		return "", "", "", err
	}

	if exists {
		return s.retrieveExistingKeys(userID, network)
	}

	return s.generateAndSaveKeys(userID, network)
}

func (s *KeyService) retrieveExistingKeys(userID int, network string) (string, string, string, error) {
	address, publicKey, encryptedPrivateKey, err := s.repository.GetKey(userID, network)
	if err != nil {
		return "", "", "", err
	}

	privateKey, err := encryption.Decrypt(encryptedPrivateKey)
	if err != nil {
		return "", "", "", err
	}

	return address, publicKey, privateKey, nil
}

func (s *KeyService) generateAndSaveKeys(userID int, network string) (string, string, string, error) {
	generator, exists := s.generators[network]
	if !exists {
		return "", "", "", errors.ErrUnsupportedNetwork
	}

	address, publicKey, privateKey, err := generator.GenerateKeyPair()
	if err != nil {
		return "", "", "", err
	}

	encryptedPrivateKey, err := encryption.Encrypt(privateKey)
	if err != nil {
		return "", "", "", err
	}

	log.Printf("Generated encrypted private key: %s", encryptedPrivateKey)

	err = s.repository.SaveKey(userID, network, address, publicKey, encryptedPrivateKey)
	if err != nil {
		return "", "", "", err
	}

	return address, publicKey, privateKey, nil
}
