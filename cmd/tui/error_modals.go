package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ErrorModal returns a modal that will
// display the error
// return to 'dest' when the user exits the modal
// NOTE: The user should always set the focus or root of the program with this
// modal
func (a *App) ErrorModal(errMsg string, dest tview.Primitive) *tview.Modal {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"OK"}).
		SetButtonBackgroundColor(tcell.ColorBlack).
		SetText(errMsg).
		SetTextColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.App.SetRoot(dest, true)
		})

	modal.SetTitle(" Error! ")
	modal.SetTitleColor(tcell.ColorRed)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}

// ExitErrorModal returns a modal that will:
// display the error
// when accepted, will safely quit the program
// NOTE: The user should always set the focus or root of the program with this
// modal
func (a *App) ExitErrorModal(errMsg string) *tview.Modal {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"OK"}).
		SetButtonBackgroundColor(tcell.ColorBlack).
		SetText(errMsg).
		SetTextColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.App.Stop()
		})

	modal.SetTitle(" Error! ")
	modal.SetTitleColor(tcell.ColorRed)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}
