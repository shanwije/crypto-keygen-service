package keys

import (
	"crypto-keygen-service/internal/bitcoin"
	"crypto-keygen-service/internal/errors"
	"crypto-keygen-service/internal/ethereum"
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
