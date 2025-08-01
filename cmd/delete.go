/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific item frmo your vault",
	Long: fmt.Sprintf(`%s

'delete' deletes a specific source name from your vault. This HAS to be case
sensitive.
Ex.
	$ gopass delete google
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := DeleteCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'delete' command: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// DeleteCmdHandler is the handler function that encapsulates the delete logic
func DeleteCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New(
			"too many or not enough arugments for 'delete'. see 'help' for correct usage",
		)
	}

	itemToDelete := args[0]

	cfg, err := utils.CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("cannot access, need to login")
	}

	err = DeleteItemInVault(cfg, itemToDelete, os.Stdin)
	if err != nil {
		return err
	}

	return nil
}

// DeleteItemInVault encapsulates the logic for deleting 'name' from the vault
// if it exists. If not, it will error and print a message out to user.
func DeleteItemInVault(cfg model.Config, name string, r io.Reader) error {
	// TODO: should I add all the open and decrypt into the confirm conditional?
	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()

	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		return fmt.Errorf("nothing in your vault")
	}

	confirm, err := utils.ConfirmPrompt(utils.DeletePrompt, name, r)
	if !confirm && err != nil {
		return fmt.Errorf("failed to confirm deletion: %v", err)
	}

	if confirm {
		found := false
		for i, v := range entries {
			if name == v.Name {
				entries = slices.Delete(entries, i, i+1)
				fmt.Printf("Deleted %s from your vault\n", name)
				found = true
			}
		}

		if !found {
			return fmt.Errorf("%s not found in vault", name)
		}
	}

	b, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("delete::failed to encrypt: %v", err)
	}

	utils.WriteToFile(f, b)

	return nil
}
