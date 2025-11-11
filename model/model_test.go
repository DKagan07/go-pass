package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNonce(t *testing.T) {
	nonce := GenerateNonce()

	assert := assert.New(t)
	assert.NotNil(nonce)
	assert.Len(nonce, NONCE_SIZE)

	// Test that nonces are different each time
	nonce2 := GenerateNonce()
	assert.NotEqual(nonce, nonce2)
}
