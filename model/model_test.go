package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNonce(t *testing.T) {
	assert := assert.New(t)
	nonce, err := GenerateNonce()
	assert.NoError(err)

	assert.NotNil(nonce)
	assert.Len(nonce, NONCE_SIZE)

	// Test that nonces are different each time
	nonce2, err := GenerateNonce()
	assert.NoError(err)
	assert.NotEqual(nonce, nonce2)
}
