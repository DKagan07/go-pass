/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-pass/cmd"
	"go-pass/utils"
)

const FileExistsErr = "file exists"

func main() {
	initProgram()
	cmd.Execute()
}

// initProgram will check to see if the file for which we will be storing the
// passwords exists. If it doesn't, it'll create one
func initProgram() {
	f := utils.OpenVault()
	defer f.Close()
}
