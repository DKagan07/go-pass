package tui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-pass/model"
)

func TestFindFilteredVaultIndex(t *testing.T) {
	assert := assert.New(t)
	app, cleanup := NewTestApp(t)
	defer cleanup()

	timestamp1 := time.Now().UnixMilli()
	timestamp2 := timestamp1 + 1000
	timestamp3 := timestamp2 + 1000

	entries := []model.VaultEntry{
		{Name: "Entry1", Username: "user1", UpdatedAt: timestamp1},
		{Name: "Entry2", Username: "user2", UpdatedAt: timestamp2},
		{Name: "Entry3", Username: "user3", UpdatedAt: timestamp3},
	}

	app.Vault = entries

	t.Run("find first entry", func(t *testing.T) {
		idx := app.findFilteredVaultIndex(entries[0])
		assert.Equal(0, idx)
	})

	t.Run("find middle entry", func(t *testing.T) {
		idx := app.findFilteredVaultIndex(entries[1])
		assert.Equal(1, idx)
	})

	t.Run("find last entry", func(t *testing.T) {
		idx := app.findFilteredVaultIndex(entries[2])
		assert.Equal(2, idx)
	})

	t.Run("entry not found", func(t *testing.T) {
		nonExistentEntry := model.VaultEntry{
			Name:      "NonExistent",
			Username:  "nobody",
			UpdatedAt: timestamp3 + 9999,
		}
		idx := app.findFilteredVaultIndex(nonExistentEntry)
		assert.Equal(-1, idx)
	})

	t.Run("empty vault", func(t *testing.T) {
		emptyApp, cleanup := NewTestApp(t)
		defer cleanup()

		idx := emptyApp.findFilteredVaultIndex(entries[0])
		assert.Equal(-1, idx)
	})
}
