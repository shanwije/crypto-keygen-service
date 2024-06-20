package bitcoin

import (
	"crypto-keygen-service/internal/network_factory"
	"crypto-keygen-service/internal/util/errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"
)

type BitcoinKeyGen struct {
	MasterSeed []byte
}

func (g *BitcoinKeyGen) GenerateKeyPairAndAddress(userID int) (network_factory.KeyPairAndAddress, error) {
	// Derive a user-specific seed using HMAC-SHA256
	userSeed := deriveUserSeed(g.MasterSeed, userID)

	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), userSeed)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate Bitcoin address")
		return network_factory.KeyPairAndAddress{}, errors.NewKeyGenError(500, "Failed to generate Bitcoin address")
	}
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())
	privateKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
	if err != nil {
		logrus.WithError(err).Error("Failed to encode Bitcoin private key to WIF")
		return network_factory.KeyPairAndAddress{}, errors.NewKeyGenError(500, "Failed to encode Bitcoin private key to WIF")
	}

	logrus.Info("Generated Bitcoin key pair")

	return network_factory.KeyPairAndAddress{
		Address:    address.EncodeAddress(),
		PublicKey:  publicKeyHex,
		PrivateKey: privateKeyWIF.String(),
	}, nil
}

func deriveUserSeed(masterSeed []byte, userID int) []byte {
	h := hmac.New(sha256.New, masterSeed)
	binary.Write(h, binary.BigEndian, int64(userID))
	return h.Sum(nil)
}
