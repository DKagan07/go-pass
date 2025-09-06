package utils

import (
	"os"
	"path"
)

// cleanup is a helper function to delete the vault and config files in the cmd
// tests
func TestCleanup() {
	_ = os.Remove(path.Join(VAULT_PATH, TEST_VAULT_NAME))
	_ = os.Remove(path.Join(CONFIG_PATH, TEST_CONFIG_NAME))
}
