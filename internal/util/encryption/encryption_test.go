package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleEncryptionKey = "4GRrhM8ClnrSmCrDvyFzPKdkJF9NcRkKwxlmIrsYhx0="

func TestMain(m *testing.M) {
	err := Setup(sampleEncryptionKey)
	if err != nil {
		panic("Failed to set up encryption: " + err.Error())
	}

	m.Run()
}

func TestEncryptDecrypt(t *testing.T) {
	originalText := "cRwricgLMcpyF4nXqJS8gDdfrtpfqjPmkq9K7EdUkJf79yPffY6N"

	encryptedText, err := Encrypt(originalText)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedText)

	decryptedText, err := Decrypt(encryptedText)
	assert.NoError(t, err)
	assert.Equal(t, originalText, decryptedText)
}
