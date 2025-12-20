package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
)

func TestDeleteFromVault(t *testing.T) {
	assert := assert.New(t)
	entries := []model.VaultEntry{
		{Name: "Entry1", Username: "user1", UpdatedAt: time.Now().UnixMilli()},
		{Name: "Entry2", Username: "user2", UpdatedAt: time.Now().UnixMilli()},
		{Name: "Entry3", Username: "user3", UpdatedAt: time.Now().UnixMilli()},
	}

	app, cleanup := NewTestAppWithData(t, entries)
	defer cleanup()

	assert.Equal(3, len(app.Vault))

	app.DeleteFromVault(1)

	assert.Equal(2, len(app.Vault))
	assert.Equal("Entry1", app.Vault[0].Name)
	assert.Equal("user1", app.Vault[0].Username)
	assert.Equal("Entry3", app.Vault[1].Name)
	assert.Equal("user3", app.Vault[1].Username)
}
