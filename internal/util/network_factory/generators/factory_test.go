package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyGenerator(t *testing.T) {
	generator, err := GetKeyGenerator("bitcoin")
	assert.NoError(t, err)
	assert.NotNil(t, generator)
	generator, err = GetKeyGenerator("ethereum")
	assert.NoError(t, err)
	assert.NotNil(t, generator)
	generator, err = GetKeyGenerator("unsupported")
	assert.Error(t, err)
	assert.Nil(t, generator)
}
