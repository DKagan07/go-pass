package utils

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/testutils"
)

const TEST_FILE_NAME = "test.json"

// cleanup is a helper function to delete the test file
func cleanup() {
	fName := path.Join(VAULT_PATH, TEST_FILE_NAME)
	os.Remove(fName)
}

func TestCreateVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	f, err := CreateVault(TEST_FILE_NAME, key)
	assert.NoError(err)
	defer f.Close()

	assert.DirExists(VAULT_PATH)
	assert.NotNil(f)
	assert.FileExists(path.Join(VAULT_PATH, TEST_FILE_NAME))
}

func TestOpenVault(t *testing.T) {
	testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	assert := assert.New(t)

	key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(err)

	f, err := CreateVault(TEST_FILE_NAME, key)
	assert.NoError(err)
	f.Close()

	f2, err := OpenVault(TEST_FILE_NAME)
	assert.NoError(err)
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
		config   *model.Config
		expected bool
	}{
		{
			name: "time is before 30 mins",
			config: &model.Config{
				LastVisited: time.Now().Add(time.Minute * -12).UnixMilli(),
				Timeout:     testutils.THIRTY_MINUTES,
			},
			expected: true,
		},
		{
			name: "time is more than 30 mins",
			config: &model.Config{
				LastVisited: time.Now().Add(time.Minute * -45).UnixMilli(),
				Timeout:     testutils.THIRTY_MINUTES,
			},
			expected: false,
		},
		{
			name: "different timeout in config before 30 mins",
			config: &model.Config{
				LastVisited: time.Now().Add(time.Minute * -45).UnixMilli(),
				Timeout:     time.Hour.Milliseconds(),
			},
			expected: true,
		},
		{
			name: "different timeout in config after range",
			config: &model.Config{
				LastVisited: time.Now().Add(time.Hour * -2).UnixMilli(),
				Timeout:     time.Hour.Milliseconds(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now().UnixMilli()
			assert.Equal(t, tt.expected, IsAccessBeforeLogin(tt.config, now))
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
		t.Run(tt.name, func(t *testing.T) {
			testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
			defer testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
			assert := assert.New(t)

			key, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
			assert.NoError(err)

			if tt.configPresent {
				fmt.Println("Creating config")
				f, err := CreateConfig(
					testutils.TEST_VAULT_NAME,
					testutils.TEST_MASTER_PASSWORD,
					testutils.TEST_CONFIG_NAME,
					key,
				)
				assert.NoError(err)
				defer f.Close()
			}

			fmt.Println("Checking config")
			cfg, err := CheckConfig(testutils.TEST_CONFIG_NAME, key)

			if tt.configPresent {
				fmt.Println("Checking config if present")
				assert.NotNil(cfg)
				assert.Equal(&model.Config{
					MasterPassword: testutils.TEST_MASTER_PASSWORD,
					VaultName:      testutils.TEST_VAULT_NAME,
					LastVisited:    cfg.LastVisited,
					Timeout:        testutils.THIRTY_MINUTES,
				}, cfg)
				assert.NoError(err)
			} else {
				fmt.Println("Checking config if not present")
				assert.Error(err)
				assert.Nil(cfg)
			}
		})

		cleanup()
		time.Sleep(time.Millisecond * 100)
	}
}
