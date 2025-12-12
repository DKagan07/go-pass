// TODO: DO THIS
/*
Copyright © 2025 DKagan07
*/
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"go-pass/model"
	"go-pass/utils"
)

// viewCmd represents the view command
var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current state of your config",
	Long: `'view' views the current state of your config. This does not show
your master password for security reasons.

Ex.
	$ gopass config view
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ViewCmdHandler(cmd, args); err != nil {
			fmt.Printf("error with view: %v\n", err)
			return
		}
	},
}

func ViewCmdHandler(cmd *cobra.Command, args []string) error {
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

	PrintConfig(cfg)

	return nil
}

func PrintConfig(cfg *model.Config) {
	fmt.Print(strings.Repeat("*", 8))
	fmt.Print(" Config ")
	fmt.Print(strings.Repeat("*", 8) + "\n")

	fmt.Printf("Vault name: %s\n", cfg.VaultName)
	fmt.Printf("Master Password: ******\n")
	fmt.Printf("Timeout: %s\n", convertTimeMsToDuration(cfg.Timeout))
	fmt.Println(strings.Repeat("*", 24))
}

func convertTimeMsToDuration(ms int64) string {
	dur := time.Duration(ms) * time.Millisecond

	hrs := int(dur.Hours())
	mins := int(dur.Minutes()) % 60
	return fmt.Sprintf("%d hours %d minutes", hrs, mins)
}
