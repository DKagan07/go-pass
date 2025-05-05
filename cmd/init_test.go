package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/utils"
)

func TestDoesConfigExist(t *testing.T) {
	tests := []struct {
		name             string
		doesConfigExists bool
	}{
		{
			name:             "config does not exist",
			doesConfigExists: false,
		},
		{
			name:             "config exists",
			doesConfigExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.doesConfigExists {
				f, err := utils.CreateConfig(
					TEST_VAULT_NAME,
					TEST_MASTER_PASSWORD,
					TEST_CONFIG_NAME,
				)
				assert.NoError(t, err)
				f.Close()
			}
			assert.Equal(t, tt.doesConfigExists, DoesConfigExist(TEST_CONFIG_NAME))
		})
	}
}

func TestCreateFiles(t *testing.T) {
	cleanup()
	defer cleanup()

	err := CreateFiles(TEST_VAULT_NAME, TEST_CONFIG_NAME, TEST_MASTER_PASSWORD)
	assert.NoError(t, err)
}

func TestEnsureVaultName(t *testing.T) {
	tests := []struct {
		name      string
		vaultName string
		good      bool
	}{
		{
			name:      "good vault name",
			vaultName: "good.json",
			good:      true,
		},
		{
			name:      "bad vault name",
			vaultName: "bad",
			good:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := EnsureVaultName(tt.vaultName)
			if tt.good {
				assert.Equal(t, s, tt.vaultName)
			} else {
				assert.Equal(t, fmt.Sprintf("%s.json", tt.vaultName), s)
			}
		})
	}
}
