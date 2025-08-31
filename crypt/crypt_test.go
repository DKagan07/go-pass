package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testEntry1 = "test1"
	testEntry2 = "test2"
)

func TestGetAESKey(t *testing.T) {
	key := GetAESKey()

	assert := assert.New(t)
	assert.NotNil(key)
	assert.Len(key, KEY_SIZE)
}

func TestGenerateNonce(t *testing.T) {
	nonce := GenerateNonce()

	assert := assert.New(t)
	assert.NotNil(nonce)
	assert.Len(nonce, NONCE_SIZE)

	// Test that nonces are different each time
	nonce2 := GenerateNonce()
	assert.NotEqual(nonce, nonce2)
}
