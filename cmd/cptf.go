/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var cptfCmd = &cobra.Command{
	Use:   "cptf",
	Short: "[C]reate [P]lain[T]ext [F]ile",
	Long: fmt.Sprintf(`%s


Creates a plain-text file of your vault, showing all passwords
This should only be used if using the previous way before the cutover
Use this sparringly, and make sure to clean up after using this
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := CptfCmdHandler(cmd, args); err != nil {
			fmt.Printf("Error with `cptf` command: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cptfCmd)
}

func CptfCmdHandler(cmd *cobra.Command, args []string) error {
	cfg, err := utils.CheckConfig("", &model.MasterAESKeyManager{})
	if err != nil {
		return err
	}

	vaultF, err := utils.OpenVault(cfg.VaultName)
	if err != nil {
		return err
	}
	defer vaultF.Close()

	ve := crypt.DecryptVault(vaultF, &model.MasterAESKeyManager{}, true)

	outF, err := os.OpenFile("out", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer outF.Close()

	decryptedVaultEntries := make([]model.DecryptedEntry, len(ve))

	for i, v := range ve {
		decryptedVaultEntries[i] = model.DecryptedEntry{
			Name:      v.Name,
			Username:  v.Username,
			Password:  crypt.DecryptPassword(v.Password, &model.MasterAESKeyManager{}, true),
			Notes:     v.Notes,
			UpdatedAt: v.UpdatedAt,
		}
	}

	b, err := json.Marshal(decryptedVaultEntries)
	if err != nil {
		return err
	}
	_, err = outF.Write(b)
	if err != nil {
		return err
	}

	outF.Sync()
	outF.Seek(0, io.SeekStart)

	// readBytes, err := io.ReadAll(outF)
	// if err != nil {
	// 	return err
	// }
	//
	// var entries []model.DecryptedEntry
	// if err := json.Unmarshal(readBytes, &entries); err != nil {
	// 	return err
	// }
	//
	// for _, e := range entries {
	// 	fmt.Printf(
	// 		"source:%s\nusername:%s\npassword:%s\nnotes:%s\n",
	// 		e.Name,
	// 		e.Username,
	// 		string(e.Password),
	// 		e.Notes,
	// 	)
	// }

	return nil
}
