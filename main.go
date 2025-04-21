/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-pass/cmd"
	"go-pass/utils"
)

const (
	SECRET_PASSWORD_KEY string = "SECRET_PASSWORD_KEY"
)

func main() {
	initProgram()
	cmd.Execute()
}

func initProgram() {
	// Just ensure that the file exists on startup
	f := utils.OpenVault()
	f.Close()

	// Ensure that the key exists on startup
	utils.GetAESKey()
}
