/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
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
		DeleteCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func DeleteCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New(
			"Too many or not enough arugments for 'delete'. See 'help' for correct usage.",
		)
	}

	itemToDelete := args[0]

	cfg, err := CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("Cannot access, need to login")
	}

	err = DeleteItemInVault(cfg, itemToDelete)
	if err != nil {
		return err
	}

	return nil
}

func DeleteItemInVault(cfg model.Config, name string) error {
	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()

	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		return fmt.Errorf("Nothing in your vault!")
	}

	found := false
	for i, v := range entries {
		if name == v.Name {
			entries = slices.Delete(entries, i, i+1)
			fmt.Printf("Deleted %s from your vault\n", name)
			found = true
		}
	}

	if !found {
		fmt.Printf("%s not found\n", name)
		return fmt.Errorf("%s not found in vault.", name)
	}

	b, err := crypt.EncryptVault(entries)
	if err != nil {
		log.Fatalf("delete::failed to encrypt: %v", err)
	}

	utils.WriteToFile(f, b)

	return nil
}
