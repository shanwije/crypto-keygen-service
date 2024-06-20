package ethereum

import (
	"crypto-keygen-service/internal/util/errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type EthereumKeyGen struct {
	MasterSeed []byte
}

func (g *EthereumKeyGen) GenerateKeyPair(userID int) (string, string, string, error) {
	// Derive a user-specific seed using HMAC-SHA256
	userSeed := deriveUserSeed(g.MasterSeed, userID)
	privateKey, err := crypto.ToECDSA(userSeed)
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

func deriveUserSeed(masterSeed []byte, userID int) []byte {
	h := hmac.New(sha256.New, masterSeed)
	binary.Write(h, binary.BigEndian, int64(userID))
	return h.Sum(nil)
}
