package utils

import (
	"fmt"
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

// TODO: UPDATE THESE TESTS FOR THE CREATE CONFIG AND VAULT FUNCTIONS
// This needs to happen in ALL of the command test files too
func TestCreateVault(t *testing.T) {
	cleanup()
	defer cleanup()
	assert := assert.New(t)

	f, err := CreateVault(TEST_FILE_NAME)
	assert.NoError(err)
	defer f.Close()

	assert.DirExists(VAULT_PATH)
	assert.NotNil(f)
	assert.FileExists(path.Join(VAULT_PATH, TEST_FILE_NAME))
}

func TestOpenVault(t *testing.T) {
	cleanup()
	defer cleanup()
	assert := assert.New(t)

	f, err := CreateVault(TEST_FILE_NAME)
	assert.NoError(err)
	f.Close()

	f2 := OpenVault(TEST_FILE_NAME)
	defer f2.Close()

	fPath := path.Join(VAULT_PATH, TEST_FILE_NAME)

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

func TestCheckConfig(t *testing.T) {
	cleanup()
	tests := []struct {
		name          string
		configPresent bool
	}{
		{
			name:          "config not created",
			configPresent: false,
		},
		{
			name:          "config is created",
			configPresent: true,
		},
	}

	for _, tt := range tests {
		defer cleanup()
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			if tt.configPresent {
				fmt.Println("Creating config")
				f, err := CreateConfig(
					TEST_VAULT_NAME,
					TEST_MASTER_PASSWORD,
					TEST_CONFIG_NAME,
				)
				assert.NoError(err)
				defer f.Close()
			}

			fmt.Println("Checking config")
			cfg, err := CheckConfig(TEST_CONFIG_NAME)

			time := cfg.LastVisited
			if tt.configPresent {
				fmt.Println("Checking config if present")
				assert.Equal(model.Config{
					MasterPassword: TEST_MASTER_PASSWORD,
					VaultName:      TEST_VAULT_NAME,
					LastVisited:    time,
				}, cfg)
				assert.NoError(err)
			} else {
				fmt.Println("Checking config if not present")
				assert.Error(err)
				assert.Equal(model.Config{}, cfg)
			}
		})

		time.Sleep(time.Millisecond * 100)
	}
}
