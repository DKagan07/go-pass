/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var uploadCptf = &cobra.Command{
	Use:   "upload",
	Short: "Upload Cptf output",
	Long: fmt.Sprintf(`%s

Uploads a file created from the 'cptf' command to re-populate the vault
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		if err := UploadCptfCmdHandler(cmd, args); err != nil {
			fmt.Printf("Error with `uploadCptf` command: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCptf)
}

// UploadCptfCmdHandler is a function that handles the upload command
// The main purpose of this is to upload the output of the cptf command to
// create a new vault
func UploadCptfCmdHandler(cmd *cobra.Command, args []string) error {
	f, err := os.OpenFile("out", os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("readFile: %+v", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("readall: %+v", err)
	}

	var ve []model.DecryptedEntry
	if err := json.Unmarshal(b, &ve); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	// Get masterpass -> create new keyring
	vaultName, err := utils.GetInputFromUser(os.Stdin, "Vault Name")
	vaultName = EnsureVaultName(vaultName)

	passBytes, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		return err
	}

	keyring := model.NewMasterAESKeyManager(string(passBytes))
	keyring.InitializeKeychain()

	bPassword, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypting password: %v", err)
	}

	cfgF, err := utils.CreateConfig(vaultName, bPassword, "", keyring)
	if err != nil {
		return err
	}
	defer cfgF.Close()

	now := time.Now().UnixMilli()
	cfg := model.Config{
		MasterPassword: bPassword,
		VaultName:      vaultName,
		LastVisited:    now,
		Timeout:        utils.THIRTY_MINUTES,
	}

	cfgCiphertext, err := crypt.EncryptConfig(&cfg, keyring)
	if err != nil {
		return err
	}
	utils.WriteToFile(cfgF, cfgCiphertext)

	vf, err := utils.CreateVault(vaultName, keyring)
	if err != nil {
		return err
	}
	defer vf.Close()

	v := make([]model.VaultEntry, len(ve))
	for i, dv := range ve {
		encryptedPass, err := keyring.Encrypt([]byte(dv.Password))
		if err != nil {
			return err
		}

		v[i] = model.VaultEntry{
			Name:     dv.Name,
			Username: dv.Username,
			Password: []byte(encryptedPass),
			Notes:    dv.Notes,
		}
	}

	ciphertext, err := crypt.EncryptVault(v, keyring)
	if err != nil {
		return err
	}

	utils.WriteToFile(vf, ciphertext)
	return nil
}
