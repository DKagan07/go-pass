package utils

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
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
	defer cleanup()

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

func TestIsAccessBeforeLogin(t *testing.T) {
	tests := []struct {
		name     string
		now      int64
		expected bool
	}{
		{
			name:     "time is before 30 mins",
			now:      time.Now().UnixMilli(),
			expected: true,
		},
		{
			name:     "time is after 30 mins",
			now:      time.Now().UnixMilli() - time.Hour.Milliseconds(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().UnixMilli()
			cfg := model.Config{
				MasterPassword: []byte("test"),
				VaultName:      "pass.json",
				LastVisited:    tt.now,
			}
			assert.Equal(t, tt.expected, IsAccessBeforeLogin(cfg, now))
		})
	}
}
