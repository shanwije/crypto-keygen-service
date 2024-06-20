package bitcoin

import (
	"crypto-keygen-service/internal/util/errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	log "github.com/sirupsen/logrus"
)

type BitcoinKeyGen struct {
	MasterSeed []byte
}

func (g *BitcoinKeyGen) GenerateKeyPair(userID int) (string, string, string, error) {
	// Derive a user-specific seed using HMAC-SHA256
	userSeed := deriveUserSeed(g.MasterSeed, userID)

	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), userSeed)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		log.WithError(err).Error("Failed to generate Bitcoin address")
		return "", "", "", errors.NewKeyGenError(500, "Failed to generate Bitcoin address")
	}
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())
	privateKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
	if err != nil {
		log.WithError(err).Error("Failed to encode Bitcoin private key to WIF")
		return "", "", "", errors.NewKeyGenError(500, "Failed to encode Bitcoin private key to WIF")
	}

	log.Info("Generated Bitcoin key pair")

	return address.EncodeAddress(), publicKeyHex, privateKeyWIF.String(), nil
}

func deriveUserSeed(masterSeed []byte, userID int) []byte {
	h := hmac.New(sha256.New, masterSeed)
	binary.Write(h, binary.BigEndian, int64(userID))
	return h.Sum(nil)
}
