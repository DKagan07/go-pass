package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
)

func TestFilterVaultBySearch(t *testing.T) {
	assert := assert.New(t)
	app, cleanup := NewTestApp(t)
	defer cleanup()

	app.AddToVault("Gmail", "notes1", "user1", "pass1")
	app.AddToVault("GitHub", "notes2", "user2", "pass2")
	app.AddToVault("Ubiquiti", "notes3", "user3", "pass3")
	app.AddToVault("AmEx", "notes4", "user4", "pass4")

	// Alphebetize the list
	app.PopulateVaultList()

	assert.Equal(4, len(app.Vault))
	assert.Equal("AmEx", app.Vault[0].Name)
	assert.Equal("notes4", app.Vault[0].Notes)
	assert.Equal("user4", app.Vault[0].Username)

	searchText := "g"

	filteredEntries := FilterVaultEntries(app.Vault, searchText)
	assert.Len(filteredEntries, 2)
	assert.Equal("GitHub", filteredEntries[0].Name)
	assert.Equal("notes2", filteredEntries[0].Notes)
	assert.Equal("user2", filteredEntries[0].Username)
	assert.Equal("pass2", crypt.DecryptPassword(filteredEntries[0].Password, app.Keyring, false))

	assert.Equal("Gmail", filteredEntries[1].Name)
	assert.Equal("notes1", filteredEntries[1].Notes)
	assert.Equal("user1", filteredEntries[1].Username)
	assert.Equal("pass1", crypt.DecryptPassword(filteredEntries[1].Password, app.Keyring, false))
}
