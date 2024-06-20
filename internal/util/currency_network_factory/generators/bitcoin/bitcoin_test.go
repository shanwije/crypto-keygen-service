package bitcoin_test

import (
	"crypto-keygen-service/internal/util/currency_network_factory/generators/bitcoin"
	"strings"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	masterSeed := []byte("test-master-seed-1234")
	keyGen := &bitcoin.BitcoinKeyGen{MasterSeed: masterSeed}

	userID := 1

	address1, publicKey1, privateKey1, err := keyGen.GenerateKeyPair(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	address2, publicKey2, privateKey2, err := keyGen.GenerateKeyPair(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address1 != address2 {
		t.Errorf("Expected the same address for the same user ID, got %s and %s", address1, address2)
	}

	if publicKey1 != publicKey2 {
		t.Errorf("Expected the same public key for the same user ID, got %s and %s", publicKey1, publicKey2)
	}

	if privateKey1 != privateKey2 {
		t.Errorf("Expected the same private key for the same user ID, got %s and %s", privateKey1, privateKey2)
	}

	if len(address1) == 0 {
		t.Errorf("Expected a valid Bitcoin address, got an empty string")
	} else if !strings.HasPrefix(address1, "1") {
		t.Errorf("Expected a Bitcoin legacy address starting with '1', got %s", address1)
	}

	if len(publicKey1) == 0 {
		t.Errorf("Expected a valid public key, got an empty string")
	}

	if len(privateKey1) == 0 {
		t.Errorf("Expected a valid private key, got an empty string")
	}
}
