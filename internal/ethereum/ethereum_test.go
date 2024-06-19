package ethereum_test

import (
	"crypto-keygen-service/internal/ethereum"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keyGen := &ethereum.EthereumKeyGen{}

	address, publicKey, privateKey, err := keyGen.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(address) == 0 {
		t.Errorf("Expected a valid Ethereum address, got an empty string")
	}
	if len(publicKey) == 0 {
		t.Errorf("Expected a valid public key, got an empty string")
	}
	if len(privateKey) == 0 {
		t.Errorf("Expected a valid private key, got an empty string")
	}
}
