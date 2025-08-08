/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"

	"go-pass/model"
	"go-pass/utils"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all storage of passwords",
	Long: fmt.Sprintf(`%s

'clean' removes your config and vault from your computer. This is a permanent
event and needs to be done with clear intentions. A simple 'y' or 'n' is needed
at the prompt.
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := CleanCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'clean' command: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

// CleanCmdHandler is the handler function for the 'clean' command.
func CleanCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("too many arugments for 'clean'. see 'help' for correct usage")
	}

	cfg, err := utils.CheckConfig("")
	if err != nil {
		return fmt.Errorf("config file does not exist: %v", err)
	}

	return CleanFiles(cfg, os.Stdin)
}

// Clean files separates the logic from the handler. This prompts the user and
// deletes the files if the user says yes.
func CleanFiles(cfg model.Config, r io.Reader) error {
	clean, err := utils.ConfirmPrompt(utils.CleanPrompt, "", os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to confirm clean: %v", err)
	}

	if clean {
		if err := RemoveConfig(""); err != nil {
			return fmt.Errorf("error removing config: %v", err)
		}
		if err := RemoveVault(cfg.VaultName); err != nil {
			return fmt.Errorf("error removing vault: %v", err)
		}
	}

	return nil
}

// RemoveConfig encapsulates the logic of removing the config.
func RemoveConfig(configFP string) error {
	if configFP == "" {
		configFP = utils.CONFIG_FILE
	} else {
		configFP = path.Join(utils.CONFIG_PATH, configFP)
	}

	if err := os.Remove(configFP); err != nil {
		return fmt.Errorf("error removing file: %v", err)
	}
	fmt.Println("Removed config.")
	return nil
}

// RemoveVault encapsulates the logic of removing the vault, or the place where
// the passwords are stored.
func RemoveVault(vaultName string) error {
	if err := os.Remove(path.Join(utils.VAULT_PATH, vaultName)); err != nil {
		return fmt.Errorf("error removing vault: %v", err)
	}
	fmt.Println("Removed vault.")
	return nil
}
