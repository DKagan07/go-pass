/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	vaultName, err := cmd.Flags().GetString("vault-name")
	if err != nil {
		log.Fatalf("init::failed to get flag: %v", err)
	}

	if vaultName == "" {
		vaultName = "pass.json"
	}

	password, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		log.Fatalf("init::failed to get password: %v", err)
	}

	masterPass, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("init::bcrypt gen pass: %v", err)
	}

	f := utils.CreateConfig(vaultName, masterPass)
	defer f.Close()

	// cfg := crypt.DecryptConfig(f)
}
