package bitcoin_test

import (
	"crypto-keygen-service/internal/network_factory/generators/bitcoin"
	"strings"
	"testing"
)

func TestGenerateKeyPairAndAddress(t *testing.T) {
	masterSeed := []byte("test-master-seed-1234")
	keyGen := &bitcoin.BitcoinKeyGen{MasterSeed: masterSeed}

	userID := 1

	keyPair1, err := keyGen.GenerateKeyPairAndAddress(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	keyPair2, err := keyGen.GenerateKeyPairAndAddress(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if keyPair1.Address != keyPair2.Address {
		t.Errorf("Expected the same address for the same user ID, got %s and %s", keyPair1.Address, keyPair2.Address)
	}

	if keyPair1.PublicKey != keyPair2.PublicKey {
		t.Errorf("Expected the same public key for the same user ID, got %s and %s", keyPair1.PublicKey, keyPair2.PublicKey)
	}

	if keyPair1.PrivateKey != keyPair2.PrivateKey {
		t.Errorf("Expected the same private key for the same user ID, got %s and %s", keyPair1.PrivateKey, keyPair2.PrivateKey)
	}

	if len(keyPair1.Address) == 0 {
		t.Errorf("Expected a valid Bitcoin address, got an empty string")
	} else if !strings.HasPrefix(keyPair1.Address, "1") {
		t.Errorf("Expected a Bitcoin legacy address starting with '1', got %s", keyPair1.Address)
	}

	if len(keyPair1.PublicKey) == 0 {
		t.Errorf("Expected a valid public key, got an empty string")
	}

	if len(keyPair1.PrivateKey) == 0 {
		t.Errorf("Expected a valid private key, got an empty string")
	}
}
