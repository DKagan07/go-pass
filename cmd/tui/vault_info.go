package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
)

// ModalVaultInfoByVault returns the Modal primitive to display the specific
// vault entry information
func (a *App) ModalVaultInfoByVault(ve model.VaultEntry) *tview.Modal {
	decryptedPassword, err := crypt.DecryptPassword(ve.Password, a.Keyring, false)
	if err != nil {
		modal := a.ErrorModal(err.Error(), a.Root)
		a.App.SetRoot(modal, true)
	}
	text := fmt.Sprintf(`
	Name: %s
	Username: %s
	Password: %s
	Notes: %s
	`, ve.Name, ve.Username, decryptedPassword, ve.Notes)
	modal := tview.NewModal().
		AddButtons([]string{"OK", "Copy"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle(" Vault Info ")
	modal.SetText(text)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if strings.EqualFold(buttonLabel, "Copy") {
			err := clipboard.WriteAll(decryptedPassword)
			if err != nil {
				a.ErrorModal(err.Error(), a.Root)
			}
		}
		a.App.SetRoot(a.Root, true)
		a.App.SetFocus(a.VaultList)
	})

	return modal
}

// ModalVaultInfoByIdx gets the specific vault entry information by index in the
// vault to display the information
func (a *App) ModalVaultInfoByIdx(idx int) *tview.Modal {
	entry := a.Vault[idx]
	decryptedPassword, err := crypt.DecryptPassword(entry.Password, a.Keyring, false)
	if err != nil {
		modal := a.ErrorModal(err.Error(), a.Root)
		a.App.SetRoot(modal, true)
	}
	text := fmt.Sprintf(`
	Name: %s
	Username: %s
	Password: %s
	Notes: %s
	`, entry.Name, entry.Username, decryptedPassword, entry.Notes)
	modal := tview.NewModal().
		AddButtons([]string{"OK", "Copy"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle(" Vault Info ")
	modal.SetText(text)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if strings.EqualFold(buttonLabel, "Copy") {
			err := clipboard.WriteAll(decryptedPassword)
			if err != nil {
				a.ErrorModal(err.Error(), a.Root)
			}
		}
		a.App.SetRoot(a.Root, true)
		a.App.SetFocus(a.VaultList)
	})

	return modal
}

// CopyDirectlyToClipboard is triggered by pressing 'c' in the ListView and
// directly copies the password to the clipboard without displaying the
// password
func (a *App) CopyDirectlyToClipboard(idx int) *tview.Modal {
	entry := a.Vault[idx]
	decryptedPassword, err := crypt.DecryptPassword(entry.Password, a.Keyring, false)
	if err != nil {
		modal := a.ErrorModal(err.Error(), a.Root)
		a.App.SetRoot(modal, true)
	}

	clipboard.WriteAll(decryptedPassword)

	modal := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle(" Vault Info ")
	modal.SetText("Copied to clipboard!")
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		a.App.SetRoot(a.Root, true)
		a.App.SetFocus(a.VaultList)
	})

	return modal
}
