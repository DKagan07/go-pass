/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize all files and begins use of the app",
	Long: fmt.Sprintf(`%s

'init' initializes all of the files and config that is required to run the app.
Notably, there's a flag that you can customize the name of your vault. 

Remember! The Master Password you set now is used for logging in with the 
'login' command.

***
IMPORTANT: The file type should be a json file, so your name should not have
any spaces and end with '.json'.
***

Ex.
	$ gopass init
	Master Password: <insert master password here>

Ex 2.
	$ gopass init --vault-name <random_name>.json
	Master Password: <insert master password here>
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		InitCmdHandler(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().
		StringP("vault-name", "v", "", "The name of the vault file that's not the default")
}

// InitCmdHandler is the handler funciton that encapsulates the logic for
// initializing the program
func InitCmdHandler(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		fmt.Println("No arguments needed. See 'help' for more information")
		return fmt.Errorf("No arguments needed. See 'help' for more information")
	}

	if DoesConfigExist("") {
		fmt.Println("Cannot init")
		return fmt.Errorf("Cannot run this command")
	}

	vaultName, err := cmd.Flags().GetString("vault-name")
	if err != nil {
		return fmt.Errorf("init::failed to get flag: %v", err)
	}

	if vaultName == "" {
		vaultName = "pass.json"
	}

	vaultName = EnsureVaultName(vaultName)

	password, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		log.Fatalf("init::failed to get password: %v", err)
	}

	bPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypting password: %v", err)
	}

	err = CreateFiles(vaultName, "", bPassword)
	if err != nil {
		return fmt.Errorf("failed creating files: %v", err)
	}

	return nil
}

// DoesConfigExist is a helper function that returns a bool whether or not the
// config file exists; true if it does, false if it doesnt
func DoesConfigExist(cfgName string) bool {
	var cfgPath string
	if cfgName == "" {
		cfgPath = utils.CONFIG_FILE
	} else {
		cfgPath = path.Join(utils.CONFIG_PATH, cfgName)
	}
	cf, err := os.Stat(cfgPath)
	if cf != nil || os.IsExist(err) {
		return true
	}
	return false
}

// ensureVaultName ensures that the vaultName is of a .json variety. If not, it
// will add it in. Should probably make this more robust
func EnsureVaultName(s string) string {
	if strings.Contains(s, ".json") {
		return s
	}
	return fmt.Sprintf("%s.json", s)
}

// CreateFiles encapsulates the logic of creating the config and vaults, and
// closing the files
func CreateFiles(vaultName string, cfgName string, pass []byte) error {
	f, err := utils.CreateConfig(vaultName, pass, cfgName)
	if err != nil {
		return err
	}
	f.Close()

	vf, err := utils.CreateVault(vaultName)
	if err != nil {
		return err
	}
	vf.Close()

	return nil
}
