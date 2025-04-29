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

func TestCreateVault(t *testing.T) {
	cleanup()
	defer cleanup()

	f := CreateVault(TEST_FILE_NAME)
	defer f.Close()

	assert := assert.New(t)
	assert.DirExists(VAULT_PATH)
	assert.NotNil(f)
	assert.FileExists(path.Join(VAULT_PATH, TEST_FILE_NAME))
}

func TestOpenVault(t *testing.T) {
	cleanup()
	// defer cleanup()

	f := CreateVault(TEST_FILE_NAME)
	f.Close()

	f2 := OpenVault(TEST_FILE_NAME)
	defer f2.Close()

	fPath := path.Join(VAULT_PATH, TEST_FILE_NAME)

	assert := assert.New(t)
	assert.NotNil(f2)
	assert.FileExists(fPath)

	stat, err := os.Stat(fPath)
	assert.NoError(err)
	assert.NotNil(stat)
	assert.NotZero(stat.Size())
}
