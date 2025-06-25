package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
	"go-pass/utils"
)

func TestUpdateEntry(t *testing.T) {
	cleanup()
	defer cleanup()

	cF, err := utils.CreateConfig(utils.TEST_VAULT_NAME, utils.TEST_MASTER_PASSWORD, utils.TEST_CONFIG_NAME)
	assert.NoError(t, err)
	cF.Close()

	vF, err := utils.CreateVault(utils.TEST_VAULT_NAME)
	assert.NoError(t, err)
	vF.Close()

	cfg := model.Config{
		VaultName:      utils.TEST_VAULT_NAME,
		MasterPassword: utils.TEST_MASTER_PASSWORD,
		LastVisited:    time.Now().UnixMilli(),
	}

	err1 := AddToVault(vaultEntry1, model.UserInput{
		Username: vaultEntry1,
		Password: []byte(vaultEntry1),
	}, cfg, time.Now().UnixMilli())
	assert.NoError(t, err1)

	i := Inputs{
		Source:   true,
		Username: false,
		Password: false,
		Notes:    false,
	}

	is := InputSources{
		Source:   strings.NewReader("newSource\n"),
		Username: strings.NewReader("newUsername\n"),
		Notes:    strings.NewReader("newNotes\n"),
	}
	err = UpdateEntry(i, cfg, vaultEntry1, is)
	assert.NoError(t, err)

	err = PrintList("newSource", cfg)
	assert.NoError(t, err)

	err = PrintList(vaultEntry1, cfg)
	assert.Error(t, err)
}

func TestUpdateVaultEntry(t *testing.T) {
	now := time.Now().UnixMilli()
	ve := model.VaultEntry{
		Name:      vaultEntry1,
		Username:  vaultEntry1,
		Password:  []byte(vaultEntry1),
		Notes:     "",
		UpdatedAt: now,
	}

	i := Inputs{
		Source:   true,
		Username: false,
		Password: false,
		Notes:    false,
	}

	is := InputSources{
		Source: strings.NewReader("newSource\n"),
	}

	newNow := time.Now().UnixMilli()
	expected := model.VaultEntry{
		Name:      "newSource",
		Username:  vaultEntry1,
		Password:  []byte(vaultEntry1),
		Notes:     "",
		UpdatedAt: newNow,
	}

	updatedVe, err := UpdateVaultEntry(ve, i, is)
	assert.NoError(t, err)
	assert.Equal(t, expected, updatedVe)
}
