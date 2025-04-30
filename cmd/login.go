/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"go-pass/crypt"
	"go-pass/utils"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		loginCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loginCmdFunc(cmd *cobra.Command, args []string) {
	cfgFile, ok, err := utils.OpenConfig()
	if ok && err == nil {
		fmt.Println("A file is not found. Need to init.")
		return
	}
	cfg := crypt.DecryptConfig(cfgFile)

	pass, err := utils.GetPasswordFromUser(true, os.Stdin)
	if err != nil {
		log.Fatalf("login::failed reading pword: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword(cfg.MasterPassword, pass); err != nil {
		fmt.Println("Login failed")
		return
	}
	fmt.Println("Success")

	now := time.Now().UnixMilli()
	cfg.LastVisited = now

	cipherText, err := crypt.EncryptConfig(cfg)
	if err != nil {
		log.Fatalf("login::failed creating ciphertext: %v", err)
	}

	utils.WriteToFile(cfgFile, cipherText)
}
