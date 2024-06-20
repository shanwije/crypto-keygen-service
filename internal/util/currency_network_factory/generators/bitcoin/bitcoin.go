package bitcoin

import (
	"crypto-keygen-service/internal/util/errors"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	log "github.com/sirupsen/logrus"
)

type BitcoinKeyGen struct{}

func (g *BitcoinKeyGen) GenerateKeyPair() (string, string, string, error) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		log.WithError(err).Error("Failed to generate Bitcoin private key")
		return "", "", "", errors.NewAPIError(500, "Failed to generate Bitcoin private key")
	}
	publicKey := privateKey.PubKey()
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	address, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.TestNet3Params)
	if err != nil {
		log.WithError(err).Error("Failed to generate Bitcoin address")
		return "", "", "", errors.NewAPIError(500, "Failed to generate Bitcoin address")
	}
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())
	privateKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		log.WithError(err).Error("Failed to encode Bitcoin private key to WIF")
		return "", "", "", errors.NewAPIError(500, "Failed to encode Bitcoin private key to WIF")
	}

	log.Info("Generated Bitcoin key pair")

	return address.EncodeAddress(), publicKeyHex, privateKeyWIF.String(), nil
}
