package service

import (
	"crypto-keygen-service/internal/bitcoin"
	"crypto-keygen-service/internal/errors"
	"crypto-keygen-service/internal/ethereum"
	"crypto-keygen-service/internal/keys"
)

type KeyService struct {
	generators map[string]keys.KeyGenerator
}

func NewKeyService() *KeyService {
	service := &KeyService{
		generators: make(map[string]keys.KeyGenerator),
	}
	service.RegisterGenerator("bitcoin", &bitcoin.BitcoinKeyGen{})
	service.RegisterGenerator("ethereum", &ethereum.EthereumKeyGen{})

	return service
}

func (s *KeyService) RegisterGenerator(network string, generator keys.KeyGenerator) {
	s.generators[network] = generator
}

func (s *KeyService) GenerateKeyPair(userID int, network string) (string, string, string, error) {
	println("Generating key pair for user", userID, "on network", network)
	generator, exists := s.generators[network]
	if !exists {
		return "", "", "", errors.ErrUnsupportedNetwork
	}
	return generator.GenerateKeyPair()
}
