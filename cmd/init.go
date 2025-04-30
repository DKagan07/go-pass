/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: fmt.Sprintf(`%s

'init' initializes all of the files and config that is required to run the app.
Notably, there's a flag that you can customize the name of your vault.

***
IMPORTANT: The file type should be a json file, so your name should not have
any spaces and end with '.json'.
***

Ex.
	$ gopass init
	Master Password: <insert master password here>

Ex 2.
	$ gopass init --vault-name <random_name>.json
`, LongDescriptionText),
	Run: func(cmd *cobra.Command, args []string) {
		initCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().
		StringP("vault-name", "v", "", "The name of the vault file that's not the default")
}

func initCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		log.Fatalf("init::cannot run command with any arguments")
	}

	vaultName, err := cmd.Flags().GetString("vault-name")
	if err != nil {
		log.Fatalf("init::failed to get flag: %v", err)
	}

	if vaultName == "" {
		vaultName = "pass.json"
	}

	vaultName = ensureVaultName(vaultName)

	password, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		log.Fatalf("init::failed to get password: %v", err)
	}

	masterPass, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("init::bcrypt gen pass: %v", err)
	}

	// TODO: need some checks here for init-ing after an init is already done
	// Perhaps check to see if the files exist? If so, then give a generic
	// message saying "cannot run this command" or something like that.
	// Maybe need to start returning errors and things along that lines to make
	// sure that happens?

	f := utils.CreateConfig(vaultName, masterPass)
	f.Close()

	vf := utils.CreateVault(vaultName)
	vf.Close()
}

// ensureVaultName ensures that the vaultName is of a .json variety
func ensureVaultName(s string) string {
	if strings.Contains(s, ".json") {
		return s
	}
	return fmt.Sprintf("%s.json", s)
}
