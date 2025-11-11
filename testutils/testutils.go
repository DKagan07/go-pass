package testutils

import (
	"os"
	"path"
	"time"

	"go-pass/model"
)

var (
	home, _              = os.UserHomeDir()
	VAULT_PATH           = path.Join(home, ".local", "gopass")
	CONFIG_PATH          = path.Join(home, ".config", "gopass")
	TEST_VAULT_NAME      = "test-vault.json"
	TEST_CONFIG_NAME     = "test-cfg.json"
	TEST_MASTER_PASSWORD = []byte("mastahpass")
	TEST_BACKUP_NAME     = "test-backup__%s.json"
	THIRTY_MINUTES       = time.Minute.Milliseconds() * 30
)

// TestCleanup is a helper function to delete the vault and config files in tests.
// It also cleans up test keyring entries.
func TestCleanup(masterPassword string) {
	_ = os.Remove(path.Join(VAULT_PATH, TEST_VAULT_NAME))
	_ = os.Remove(path.Join(CONFIG_PATH, TEST_CONFIG_NAME))

	// Clean up test keyring entry
	keyManager := model.NewTestMasterAESKeyManager(masterPassword)
	_ = keyManager.DeleteKeychain()
}

// InitTestKeyring initializes a test keyring with a random key.
// This should be called at the start of tests that need keyring access.
// Returns the test keyring manager.
func InitTestKeyring(masterPassword string) (*model.MasterAESKeyManager, error) {
	keyManager := model.NewTestMasterAESKeyManager(masterPassword)
	if err := keyManager.InitializeKeychain(); err != nil {
		return nil, err
	}
	return keyManager, nil
}
