/*
Copyright © 2026 DKagan07
*/
package vault

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

type ExportFormat string

const (
	OnePass   ExportFormat = "1PASSWORD"
	LastPass  ExportFormat = "LASTPASS"
	BitWarden ExportFormat = "BITWARDEN"
	NordPass  ExportFormat = "NORDPASS"
	KeePass   ExportFormat = "KEEPASS"
)

var validationMap = map[ExportFormat]bool{
	OnePass:   true,
	LastPass:  true,
	BitWarden: true,
	NordPass:  true,
	KeePass:   true,
}

// ExportCmd represents the delete command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "export your vault in a specific format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("in export command")
		if err := ExportCmdHandler(cmd, args); err != nil {
			fmt.Println("Error with 'export' command: ", err)
			return
		}
	},
}

// TODO: break this up into smaller functions, especially the csv-creation ones
func ExportCmdHandler(cmd *cobra.Command, args []string) error {
	// Get export format flag
	eFormat, err := cmd.Flags().GetString("export-format")
	if err != nil {
		return fmt.Errorf("with getting string flag: %v", err)
	}

	if !validateFormatInput(strings.ToUpper(eFormat)) {
		return errors.New("not a valid format to export")
	}

	// Get vault
	passB, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}
	keyring := model.NewMasterAESKeyManager(string(passB))

	cfg, err := utils.CheckConfig("", keyring)
	if err != nil {
		return fmt.Errorf("error checking config: %v", err)
	}

	f, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return fmt.Errorf("opening vault: %v", err)
	}
	defer f.Close()

	entries, err := crypt.DecryptVault(f, keyring, false)
	if err != nil {
		return fmt.Errorf("decrypting vault: %v", err)
	}

	if len(entries) == 0 {
		return fmt.Errorf("nothing in vault")
	}

	// create file for csv output
	out, err := os.OpenFile("out.csv", os.O_RDWR|os.O_CREATE, 0o0600)
	if err != nil {
		return err
	}
	defer out.Close()

	cw := csv.NewWriter(out)

	format := ExportFormat(strings.ToUpper(eFormat))
	fmt.Println("format:", format)
	switch format {
	case OnePass:
		for _, entry := range entries {
			ce := []string{entry.Name, "", entry.Username, string(entry.Password)}
			if err := cw.Write(ce); err != nil {
				return err
			}
		}
		cw.Flush()
	case NordPass:
	case BitWarden:
		header := []string{
			"folder",
			"favorite",
			"type",
			"name",
			"notes",
			"fields",
			"reprompt",
			"login_uri",
			"login_username",
			"login_password",
			"login_totp",
		}
		if err := cw.Write(header); err != nil {
			return err
		}

		for _, entry := range entries {
			dp, err := crypt.DecryptPassword(entry.Password, keyring, false)
			if err != nil {
				return err
			}

			ce := []string{
				"Personal",
				"",
				"login",
				"",
				entry.Name,
				entry.Notes,
				"",
				"",
				"",
				entry.Username,
				dp,
				"",
			}
			if err := cw.Write(ce); err != nil {
				return err
			}
		}
		cw.Flush()

	case LastPass:
	case KeePass:
	default:
	}

	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}

func validateFormatInput(input string) bool {
	_, ok := validationMap[ExportFormat(input)]
	return ok
}
