package keys

import (
	"crypto-keygen-service/internal/bitcoin"
	"errors"
)

func GetKeyGenerator(network string) (KeyGenerator, error) {
	switch network {
	case "bitcoin":
		return &bitcoin.BitcoinKeyGen{}, nil
	default:
		return nil, errors.ErrUnsupported
	}
}
