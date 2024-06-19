package ethereum

import (
	"crypto-keygen-service/internal/util/errors"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumKeyGen struct{}

func (g *EthereumKeyGen) GenerateKeyPair() (string, string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", errors.NewAPIError(500, "Failed to generate Ethereum private key")
	}

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	publicKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey))
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return address, publicKeyHex, privateKeyHex, nil
}
