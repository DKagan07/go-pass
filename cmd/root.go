/*
Copyright © 2025 DKagan07
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var LongDescriptionText = `GoPass is a CLI tool that help stores your passwords with security in mind.
This application encrypts, hashes, and stores passwords.
Everything is local to your computer! Nothing is stored on the internet.`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopass",
	Short: "Stores and encrypts all of your sensitive passwords",
	Long:  LongDescriptionText,
}

var isOther bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().
		BoolVarP(&isOther, "isOther", "o", false, "Use other framework instead of the default")

	if err := rootCmd.ParseFlags(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if !isOther {
		err := rootCmd.Execute()
		if err != nil {
			os.Exit(1)
		}
	} else {
		TviewRun()
	}
}

func init() {}
