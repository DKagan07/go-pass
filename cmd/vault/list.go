/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the sources of your login infos",
	Long: `'list' lists all the sources of login info that's currently in your vault. A
source is, for example, 'Google', when it comes to what username and password
are stored with it.
Ex.
	$ gopass vault list
	Google
	Github
	...(etc.)
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ListCmdHandler(cmd, args); err != nil {
			fmt.Printf("Error with 'list' command: %v\n", err)
			return
		}
	},
}

// ListCmdHandler is the handler function that encapsulates the PrintList
func ListCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("no arguments needed for 'list'. see 'help' for more guidance")
	}

	bp, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(bp))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	sourceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return fmt.Errorf("list::getting string from flag: %v", err)
	}

	backups, err := cmd.Flags().GetBool("backups")
	if err != nil {
		return fmt.Errorf("list::getting bool from flag: %v", err)
	}

	if !backups {
		err = PrintList(sourceName, cfg, keyring)
		if err != nil {
			return fmt.Errorf("error printing list: %v", err)
		}
	} else {
		if err = PrintBackups(); err != nil {
			return fmt.Errorf("error printing backups: %v", err)
		}
	}

	return nil
}

// PrintList is the function that prints the list of sources in the vault.
// If a source name is provided, it will check if the source exists in the vault.
// If a source name is not provided, it will print all sources in the vault.
func PrintList(sourceName string, cfg *model.Config, key *model.MasterAESKeyManager) error {
	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()

	entries, err := crypt.DecryptVault(f, key, false)
	if err != nil {
		return fmt.Errorf("decrypting vault: %v", err)
	}

	if len(entries) == 0 {
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

	encryptedCipherText, err := crypt.EncryptVault(entries, key)
	if err != nil {
		return fmt.Errorf("list::obtaining ciphertext: %v", err)
	}

	utils.WriteToFile(f, encryptedCipherText)
	return nil
}

func PrintBackups() error {
	dirEntries, err := os.ReadDir(utils.BACKUP_DIR)
	if err != nil {
		return err
	}

	if len(dirEntries) == 0 {
		fmt.Println("No backups found")
		return nil
	}

	for _, v := range dirEntries {
		fmt.Printf("%s\n", v.Name())
	}

	return nil
}
