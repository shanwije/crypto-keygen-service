package bitcoin

import (
	"crypto-keygen-service/internal/errors"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type BitcoinKeyGen struct{}

func (g *BitcoinKeyGen) GenerateKeyPair() (string, string, string, error) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", "", errors.NewAPIError(500, "Failed to generate Bitcoin private key")
	}
	publicKey := privateKey.PubKey()
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())

	// considering this is an assignment haven't set a flag to switch testnet params
	address, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.TestNet3Params)
	if err != nil {
		return "", "", "", errors.NewAPIError(500, "Failed to generate Bitcoin address")
	}
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())
	privateKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.TestNet3Params, true)

	return address.EncodeAddress(), publicKeyHex, privateKeyWIF.String(), nil
}
