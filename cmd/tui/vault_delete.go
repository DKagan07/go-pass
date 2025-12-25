package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DeleteVaultModal returns the Modal primitive to delete a vault entry. This
// also has built-in confirmation to ensure deletion of the entry
func (a *App) DeleteVaultModal(i int) *tview.Modal {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Yes", "No"}).
		SetButtonBackgroundColor(tcell.Color103).
		SetText(fmt.Sprintf("Are you sure you want to delete %s?", a.Vault[i].Name)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if strings.EqualFold(buttonLabel, "Yes") { // This is the validation
				a.DeleteFromVault(i)

				a.PopulateVaultList()
				a.RefreshRoot()
				a.App.SetRoot(a.Root, true)
				return
			}
			a.App.SetRoot(a.Root, true)
		})

	modal.SetTitle(" Delete Vault ")
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}

// DeleteFromVault contains the business logic of removing the vault entry from
// the vault, and saves the new vault to disk
func (a *App) DeleteFromVault(vaultIdx int) {
	a.Vault = slices.Delete(a.Vault, vaultIdx, vaultIdx+1)
	a.SaveVault()
}
