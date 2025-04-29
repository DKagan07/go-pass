/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-pass/cmd"
	"go-pass/crypt"
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
	f := utils.CreateVault("")
	f.Close()

	// Ensure that the key exists on startup
	crypt.GetAESKey()
}
