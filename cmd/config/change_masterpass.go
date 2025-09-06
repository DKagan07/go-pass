/*
Copyright © 2025 DKagan07
*/
package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// changeMasterpassCmd represents the changeMasterpass command
var ChangeMasterpassCmd = &cobra.Command{
	Use:   "change_masterpass",
	Short: "Changes your master password, used for logging in",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("changeMasterpass called")
	},
}

func init() {
	// ConfigCmd.AddCommand(changeMasterpassCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changeMasterpassCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changeMasterpassCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
