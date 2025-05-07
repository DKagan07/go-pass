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
	"strings"

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
		CleanCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

// CleanCmdHandler is the handler function for the 'clean' command.
func CleanCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("Too many arugments for 'clean'. See 'help' for correct usage.")
	}

	cfg, err := CheckConfig("")
	if err != nil {
		return fmt.Errorf("Config file does not exist: %v", err)
	}

	return CleanFiles(cfg, os.Stdin)
}

// Clean files separates the logic from the handler. This prompts the user and
// deletes the files if the user says yes.
func CleanFiles(cfg model.Config, r io.Reader) error {
	fmt.Println("'clean' will remove your config and your vault.")
	ans, err := utils.GetInputFromUser(r, "Are you sure? (y/n)")
	if err != nil {
		return fmt.Errorf("Clean: Error with user input: %v", err)
	}

	if strings.EqualFold(ans, "y") {
		if err := RemoveConfig(""); err != nil {
			fmt.Println("Error removing config")
			return err
		}
		if err := RemoveVault(cfg.VaultName); err != nil {
			fmt.Println("Error removing vault")
			return err
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
		return fmt.Errorf("Error removing file: %v", err)
	}
	fmt.Println("Removed config.")
	return nil
}

// RemoveVault encapsulates the logic of removing the vault, or the place where
// the passwords are stored.
func RemoveVault(vaultName string) error {
	if err := os.Remove(path.Join(utils.VAULT_PATH, vaultName)); err != nil {
		return fmt.Errorf("Error removing vault: %v", err)
	}
	fmt.Println("Removed vault.")
	return nil
}
