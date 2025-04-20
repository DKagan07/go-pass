package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tea.Model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there is an error! %v", err)
		os.Exit(1)
	}
}
