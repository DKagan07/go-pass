package utils

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_FILE_NAME = "test.json"

// cleanup is a helper function to delete the test file
func cleanup() {
	fName := path.Join(VAULT_PATH, TEST_FILE_NAME)
	os.Remove(fName)
}

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

func TestCreateVault(t *testing.T) {
	cleanup()
	defer cleanup()

	f := CreateVault(TEST_FILE_NAME)

	assert := assert.New(t)
	assert.DirExists(VAULT_PATH)
	assert.NotNil(f)
	assert.FileExists(path.Join(VAULT_PATH, TEST_FILE_NAME))
}

// func TestOpenVault(t *testing.T) J {
// 	cleanup()
// 	defer cleanup()
//
// 	f := CreateVault(TEST_FILE_NAME)
// 	f.Close()
//
// 	f2 := OpenVault()
// }
