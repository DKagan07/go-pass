package tui

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/testutils"
	"go-pass/utils"
)

// NewTestApp creates an isolated App instance for testing
// This uses test vault/config files and test keyring
func NewTestApp(t *testing.T) (*App, func()) {
	// Initialize test keyring
	keyManager, err := testutils.InitTestKeyring(string(testutils.TEST_MASTER_PASSWORD))
	assert.NoError(t, err, "Failed to initialize test keyring")

	// Create test config
	cfgFile, err := utils.CreateConfig(
		testutils.TEST_VAULT_NAME,
		testutils.TEST_MASTER_PASSWORD,
		testutils.TEST_CONFIG_NAME,
		keyManager,
	)
	assert.NoError(t, err)

	// Create test vault
	vaultFile, err := utils.CreateVault(testutils.TEST_VAULT_NAME, keyManager)
	assert.NoError(t, err)

	// Load config
	cfg, err := crypt.DecryptConfig(cfgFile, keyManager, false)
	assert.NoError(t, err)

	// Load vault
	// vault := crypt.DecryptVault(vaultFile, keyManager, false)

	// Create App instance
	testApp := &App{
		App:           tview.NewApplication(),
		VaultFile:     vaultFile,
		Vault:         make([]model.VaultEntry, 0),
		FilteredVault: make([]model.VaultEntry, 0),
		Cfg:           cfg,
		Keyring:       keyManager,
		NumRetries:    0,
	}

	// Return cleanup function
	cleanup := func() {
		vaultFile.Close()
		cfgFile.Close()
		testutils.TestCleanup(string(testutils.TEST_MASTER_PASSWORD))
	}

	return testApp, cleanup
}

func NewTestAppWithData(t *testing.T, entries []model.VaultEntry) (*App, func()) {
	testApp, cleanup := NewTestApp(t)

	testApp.Vault = entries
	testApp.FilteredVault = entries
	testApp.SaveVault()

	return testApp, cleanup
}
