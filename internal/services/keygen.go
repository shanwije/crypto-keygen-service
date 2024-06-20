package services

import (
	"crypto-keygen-service/internal/repositories"
	"crypto-keygen-service/internal/util/currency_network_factory"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/ethereum"
	"crypto-keygen-service/internal/util/encryption"
	"crypto-keygen-service/internal/util/errors"
)

type KeyGenService struct {
	generators map[string]currency_network_factory.KeyGenerator
	repository repositories.KeyGenRepository
}

func NewKeyGenService(repo repositories.KeyGenRepository) *KeyGenService {
	service := &KeyGenService{
		generators: make(map[string]currency_network_factory.KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{})
	// Add more networks here
	return service
}

func (s *KeyGenService) RegisterGenerator(network string, generator currency_network_factory.KeyGenerator) {
	s.generators[network] = generator
}

// GetKeysAndAddress retrieves or generates keys and address for a user on a specific network.
func (s *KeyGenService) GetKeysAndAddress(userID int, network string) (string, string, string, error) {
	// Check if keys already exist in the repositories
	exists, err := s.repository.KeyExists(userID, network)
	if err != nil {
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
		return "", "", "", err
	}

	privateKey, err := encryption.Decrypt(encryptedPrivateKey)
	if err != nil {
		return "", "", "", err
	}

	return address, publicKey, privateKey, nil
}

func (s *KeyGenService) generateAndSaveKeys(userID int, network string) (string, string, string, error) {
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

	err = s.repository.SaveKey(userID, network, address, publicKey, encryptedPrivateKey)
	if err != nil {
		return "", "", "", err
	}

	return address, publicKey, privateKey, nil
}
