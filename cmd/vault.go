/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"go-pass/cmd/vault"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Parent command for vault commands",
	Long: fmt.Sprintf(`%s

'vault' is acting as a parent command for vault-related commands.
Please use the 'help' flag for additional information.
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)

	vaultCmd.AddCommand(vault.AddCmd)
	vaultCmd.AddCommand(vault.BackupCmd)
	vaultCmd.AddCommand(vault.DeleteCmd)
	vaultCmd.AddCommand(vault.GenerateCmd)
	vaultCmd.AddCommand(vault.GetCmd)
	vaultCmd.AddCommand(vault.ListCmd)
	vaultCmd.AddCommand(vault.RestoreCmd)
	vaultCmd.AddCommand(vault.SearchCmd)
	vaultCmd.AddCommand(vault.UpdateCmd)

	initVaultFlags()
}

func initVaultFlags() {
	// Get Command
	vault.GetCmd.Flags().BoolP("copy", "y", false, "Add password to clipboard, does not display information")

	// Generate Command
	specialCharsStr := "List the special characters you want to add to your password generation. If adjustment is necessary, list all the special characters you want. IMPORTANT: BE SURE TO USE SINGLE QUOTES."
	vault.GenerateCmd.Flags().IntP("length", "l", 24, "Decides length of new password")
	vault.GenerateCmd.Flags().
		StringP("add", "a", "", "Add a newly generated password to your vault")
	vault.GenerateCmd.Flags().StringP("specialChars", "c", vault.DefaultChars, specialCharsStr)

	// List Command
	vault.ListCmd.Flags().StringP("name", "n", "", "Searches your list for the specific source")
	vault.ListCmd.Flags().BoolP("backups", "b", false, "Lists your backups")

	// Update Command
	vault.UpdateCmd.Flags().BoolP("source", "s", false, "Update the source name")
	vault.UpdateCmd.Flags().BoolP("username", "u", false, "Update the login username")
	vault.UpdateCmd.Flags().BoolP("password", "p", false, "Update the password")
	vault.UpdateCmd.Flags().BoolP("notes", "t", false, "Update the notes")
}
