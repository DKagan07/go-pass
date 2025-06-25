/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the sources of your login infos",
	Long: fmt.Sprintf(`%s

'list' lists all the sources of login info that's currently in your vault. A
source is, for example, 'Google', when it comes to what username and password
are stored with it.
Ex.
	$ gopass list
	Google
	Github
	...(etc.)
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := ListCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with list")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("name", "n", "", "Searches your list for the specific source")
}

// ListCmdHandler is the handler function that encapsulates the PrintList
func ListCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Println("No arguments needed for 'list'. See 'help' for more guidance")
		return fmt.Errorf("no arguments needed for 'list'. see 'help' for more guidance")
	}

	sourceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return fmt.Errorf("list::getting string from flag: %v", err)
	}

	cfg, err := utils.CheckConfig("")
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		fmt.Println("Cannot access, need to login")
		return fmt.Errorf("cannot access, need to login")
	}

	err = PrintList(sourceName, cfg)
	if err != nil {
		return fmt.Errorf("error printing list: %v", err)
	}

	return nil
}

// PrintList is the function that prints the list of sources in the vault.
// If a source name is provided, it will check if the source exists in the vault.
// If a source name is not provided, it will print all sources in the vault.
func PrintList(sourceName string, cfg model.Config) error {
	f := utils.OpenVault(cfg.VaultName)
	defer f.Close()
	entries := crypt.DecryptVault(f)

	if len(entries) == 0 {
		fmt.Println("Nothing in your vault!")
		return fmt.Errorf("nothing in vault")
	}

	// Alphabetize the entries by `Name`
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	if sourceName == "" {
		fmt.Println("Entries:")
		for _, v := range entries {
			fmt.Printf("\t%s\n", v.Name)
		}
	} else {
		for _, v := range entries {
			if strings.EqualFold(v.Name, sourceName) {
				fmt.Printf("Yes, %s exists\n", v.Name)
				return nil
			}
		}
		return fmt.Errorf("%s does not exist", sourceName)
	}

	encryptedCipherText, err := crypt.EncryptVault(entries)
	if err != nil {
		fmt.Println("Error with 'add' command")
		return fmt.Errorf("add::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return nil
}
