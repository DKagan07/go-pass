package crypt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/testutils"
)

var (
	testEntry1 = "test1"
	testEntry2 = "test2"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err)

	passwords := []string{
		"simple123",
		"with spaces and special !@#$",
		"unicode: 你好🔐",
		strings.Repeat("a", 10000), // Long password
	}

	for _, pwd := range passwords {
		p := []byte(pwd)
		encrypted, err := EncryptPassword(p, key)
		assert.NoError(t, err)

		decrypted, err := DecryptPassword([]byte(encrypted), key, false)
		assert.NoError(t, err)

		if decrypted != pwd {
			t.Errorf("Mismatch! Original: %s, Decrypted: %s", pwd, decrypted)
		}
	}
}
