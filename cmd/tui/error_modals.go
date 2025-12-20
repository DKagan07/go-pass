package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
