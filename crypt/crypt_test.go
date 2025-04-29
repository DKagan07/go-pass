package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}
