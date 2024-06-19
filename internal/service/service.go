package service

import (
	"crypto-keygen-service/internal/bitcoin"
	"crypto-keygen-service/internal/crypto"
	"crypto-keygen-service/internal/errors"
	"crypto-keygen-service/internal/ethereum"
	"crypto-keygen-service/internal/keys"
	"crypto-keygen-service/internal/repository"
	"log"
)

type KeyService struct {
	generators map[string]keys.KeyGenerator
	repository repository.Repository
}

func NewKeyService(repo repository.Repository) *KeyService {
	service := &KeyService{
		generators: make(map[string]keys.KeyGenerator),
		repository: repo,
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{})
	// Add more networks here
	return service
}

func (s *KeyService) RegisterGenerator(network string, generator keys.KeyGenerator) {
	s.generators[network] = generator
}

func (s *KeyService) GetKeysAndAddress(userID int, network string) (string, string, string, error) {
	exists, err := s.repository.KeyExists(userID, network)
	if err != nil {
		return "", "", "", err
	}

	if exists {
		address, publicKey, encryptedPrivateKey, err := s.repository.GetKey(userID, network)
		if err != nil {
			return "", "", "", err
		}
		privateKey, err := crypto.Decrypt(encryptedPrivateKey)
		if err != nil {
			return "", "", "", err
		}
		return address, publicKey, privateKey, nil
	}

	generator, exists := s.generators[network]
	if !exists {
		return "", "", "", errors.ErrUnsupportedNetwork
	}
	address, publicKey, privateKey, err := generator.GenerateKeyPair()
	if err != nil {
		return "", "", "", err
	}
	encryptedPrivateKey, err := crypto.Encrypt(privateKey)
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
