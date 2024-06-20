package ethereum

import (
	"crypto-keygen-service/internal/util/errors"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type EthereumKeyGen struct{}

func (g *EthereumKeyGen) GenerateKeyPair() (string, string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.WithError(err).Error("Failed to generate Ethereum private key")
		return "", "", "", errors.NewAPIError(500, "Failed to generate Ethereum private key")
	}

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	publicKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey))
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	log.WithFields(log.Fields{
		"address":    address,
		"public_key": publicKeyHex,
	}).Info("Generated Ethereum key pair")

	return address, publicKeyHex, privateKeyHex, nil
}
