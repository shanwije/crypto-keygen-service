package bitcoin_test

import (
	"crypto-keygen-service/internal/util/currency_factory/generators/bitcoin"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keyGen := &bitcoin.BitcoinKeyGen{}

	address, publicKey, privateKey, err := keyGen.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(address) == 0 {
		t.Errorf("Expected a valid Bitcoin address, got an empty string")
	}
	if len(publicKey) == 0 {
		t.Errorf("Expected a valid public key, got an empty string")
	}
	if len(privateKey) == 0 {
		t.Errorf("Expected a valid private key, got an empty string")
	}
}
