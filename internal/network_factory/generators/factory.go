package generators

import (
	"crypto-keygen-service/internal/network_factory"
	"crypto-keygen-service/internal/network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/network_factory/generators/ethereum"
	"crypto-keygen-service/internal/util/errors"
)

func GetKeyGenerator(network string) (network_factory.KeyGenerator, error) {
	switch network {
	case "bitcoin":
		return &bitcoin.BitcoinKeyGen{}, nil
	case "ethereum":
		return &ethereum.EthereumKeyGen{}, nil
	default:
		return nil, errors.ErrUnsupportedNetwork
	}
}
