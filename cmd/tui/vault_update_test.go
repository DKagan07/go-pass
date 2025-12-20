package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
)

func TestValidateUpdateInputs(t *testing.T) {
	assert := assert.New(t)
	app, cleanup := NewTestApp(t)
	defer cleanup()

	app.AddToVault("Entry1", "notes1", "user1", "pass1")

	assert.Len(app.Vault, 1)

	newVaultEntry, err := app.ValidateUpdateInputs("NewEntry1", "NewUser1", "newPass1", "newNotes1")
	assert.NoError(err)
	assert.NotEqual(newVaultEntry.Name, app.Vault[0].Name)
	assert.NotEqual(newVaultEntry.Username, app.Vault[0].Username)
	assert.NotEqual(newVaultEntry.Password, app.Vault[0].Password)
	assert.NotEqual(newVaultEntry.Notes, app.Vault[0].Notes)
}

func TestUpdateVaultEntry(t *testing.T) {
	assert := assert.New(t)
	app, cleanup := NewTestApp(t)
	defer cleanup()

	app.AddToVault("Entry1", "notes1", "user1", "pass1")
	assert.Len(app.Vault, 1)
	assert.Equal("Entry1", app.Vault[0].Name)

	newVaultEntry, err := app.ValidateUpdateInputs("NewEntry1", "NewUser1", "newPass1", "newNotes1")
	assert.NoError(err)

	app.UpdateVaultEntry(0, *newVaultEntry)

	assert.Len(app.Vault, 1)
	assert.Equal("NewEntry1", app.Vault[0].Name)
	assert.NotEqual("Entry1", app.Vault[0].Name)
	assert.Equal("NewUser1", app.Vault[0].Username)
	assert.NotEqual("user1", app.Vault[0].Username)
	assert.Equal("newPass1", crypt.DecryptPassword(app.Vault[0].Password, app.Keyring, false))
	assert.NotEqual("pass1", crypt.DecryptPassword(app.Vault[0].Password, app.Keyring, false))
	assert.Equal("newNotes1", app.Vault[0].Notes)
	assert.NotEqual("notes1", app.Vault[0].Notes)
}
