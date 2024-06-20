package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/chacha20poly1305"
)

var key []byte

func Setup(keyString string) error {
	if keyString == "" {
		log.Error("ENCRYPTION_KEY not set")
		return errors.New("ENCRYPTION_KEY not set")
	}

	var err error
	key, err = base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		log.WithError(err).Error("Invalid ENCRYPTION_KEY")
		return fmt.Errorf("invalid ENCRYPTION_KEY: %w", err)
	}

	if len(key) != chacha20poly1305.KeySize {
		log.Errorf("ENCRYPTION_KEY must be %d bytes long", chacha20poly1305.KeySize)
		return fmt.Errorf("ENCRYPTION_KEY must be %d bytes long", chacha20poly1305.KeySize)
	}

	log.Info("Encryption setup successful")
	return nil
}

func Encrypt(plaintext string) (string, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		log.WithError(err).Error("Failed to create AEAD instance for encryption")
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.WithError(err).Error("Failed to generate nonce for encryption")
		return "", err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	encryptedText := base64.StdEncoding.EncodeToString(ciphertext)
	log.Debug("Encryption successful")
	return encryptedText, nil
}

func Decrypt(ciphertext string) (string, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		log.WithError(err).Error("Failed to create AEAD instance for decryption")
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		log.WithError(err).Error("Failed to decode base64 ciphertext")
		return "", err
	}

	if len(data) < aead.NonceSize() {
		log.Error("Ciphertext too short")
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:aead.NonceSize()], data[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt ciphertext")
		return "", err
	}

	log.Debug("Decryption successful")
	return string(plaintext), nil
}
