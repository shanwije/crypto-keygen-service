package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	originalText := "cRwricgLMcpyF4nXqJS8gDdfrtpfqjPmkq9K7EdUkJf79yPffY6N"

	encryptedText, err := Encrypt(originalText)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedText)

	decryptedText, err := Decrypt(encryptedText)
	assert.NoError(t, err)
	assert.Equal(t, originalText, decryptedText)
}
