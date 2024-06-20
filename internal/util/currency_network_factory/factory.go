package currency_network_factory

import (
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"crypto-keygen-service/internal/util/currency_network_factory/generators/ethereum"
	"crypto-keygen-service/internal/util/errors"
)

func GetKeyGenerator(network string) (KeyGenerator, error) {
	switch network {
	case "bitcoin":
		return &bitcoin.BitcoinKeyGen{}, nil
	case "ethereum":
		return &ethereum.EthereumKeyGen{}, nil
	default:
		return nil, errors.ErrUnsupportedNetwork
	}
}
