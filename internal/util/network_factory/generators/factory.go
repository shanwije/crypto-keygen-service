package generators

import (
	"crypto-keygen-service/internal/util/errors"
	"crypto-keygen-service/internal/util/network_factory"
	"crypto-keygen-service/internal/util/network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/network_factory/generators/ethereum"
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
