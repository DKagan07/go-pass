/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A test command to understand Cobra a bit more",
	Long:  LongDescriptionText,
	Run: func(cmd *cobra.Command, args []string) {
		helloCmdFunc(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	helloCmd.Flags().StringP("name", "n", "", "Name to greet")
}

func helloCmdFunc(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		name = "World"
	}

	fmt.Printf("Hello, %s!\n", name)
}
