package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyGenerator(t *testing.T) {
	generator, err := GetKeyGenerator("unsupported")
	assert.Error(t, err)
	assert.Nil(t, generator)
}
