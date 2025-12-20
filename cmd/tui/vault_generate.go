package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) GeneratedPasswordModal(generatedPass string) *tview.Modal {
	modal := tview.NewModal().
		AddButtons([]string{"OK", "Copy"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle(" Generated Password ")
	modal.SetText(generatedPass)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if strings.EqualFold(buttonLabel, "Copy") {
			err := clipboard.WriteAll(generatedPass)
			if err != nil {
				modal := a.ErrorModal(fmt.Sprintf("Failed to copy password: %v", err), a.Root)
				a.App.SetRoot(modal, true)
			}
		} else {
			a.App.SetRoot(a.Root, true)
		}
	})

	return modal
}
