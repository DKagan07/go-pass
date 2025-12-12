/*
Copyright © 2025 DKagan07
*/
package vault

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// searchCmd represents the search command
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a password in the vault",
	Long: `'search' searches your vault for a source that matches the search term. This
search is case insensitive, and will list all sources that match the search term.

Ex.
	$ gopass search git
	GitHub
	GitLab

Ex. 
	$ gopass search gitlab
	GitLab
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := SearchCmdHandler(cmd, args); err != nil {
			return
		}
	},
}

// func init() {
// 	rootCmd.AddCommand(searchCmd)
// }

// SearchCmdHandler is the handler function that encapsulates the SearchVault
// logic and runs some checks beforehand.
func SearchCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New(
			"too many or not enough arugments for 'search'. See 'help' for correct usage",
		)
	}

	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		return fmt.Errorf("cannot access, need to login")
	}

	searchTerm := strings.ToLower(args[0])
	return SearchVault(searchTerm, cfg, keyring)
}

// SearchVault is the function that searches the vault for a source that matches
// the search term. It will print out all sources that match the search term.
// This is a case insensitive search.
func SearchVault(searchTerm string, cfg *model.Config, key *model.MasterAESKeyManager) error {
	if searchTerm == "" {
		return fmt.Errorf("no search term provided")
	}

	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()

	entries := crypt.DecryptVault(f, key, false)

	if len(entries) == 0 {
		return fmt.Errorf("nothing in vault")
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	found := false
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Name), searchTerm) {
			found = true
			fmt.Println(e.Name)
		}
	}

	if !found {
		fmt.Println("No matches found.")
	}

	return nil
}
