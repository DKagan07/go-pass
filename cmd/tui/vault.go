package tui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/cmd/vault"
	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// PopulateVaultList is the main engine behind the TUI application. This
// repopulates the vault list from the Vault and FilteredVault, alphabezies the
// vault, and controls the button presses for actions within the VaultList
func (a *App) PopulateVaultList() {
	// Alphebetize the vault by name
	sort.Slice(a.Vault, func(i, j int) bool {
		return a.Vault[i].Name < a.Vault[j].Name
	})

	if a.SearchInput != nil {
		a.SyncFilteredVault()
	}

	a.VaultList = tview.NewList()
	for _, v := range a.FilteredVault {
		vault := v
		a.VaultList.AddItem(vault.Name, "", 0, func() {
			m := a.ModalVaultInfoByVault(vault)
			a.App.SetRoot(m, false)
		})
	}

	a.VaultList.SetBorder(true)
	a.VaultList.SetTitle(" Vault ")
	a.VaultList.SetBackgroundColor(tcell.ColorBlack)
	a.VaultList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'a':
			flex := a.ModalAddVault()
			a.App.SetRoot(flex, true)
		case 'd':
			currentIndex := a.VaultList.GetCurrentItem()
			if currentIndex >= 0 && currentIndex < len(a.Vault) {
				entry := a.FilteredVault[currentIndex]
				actualIdx := a.findFilteredVaultIndex(entry)

				if actualIdx >= 0 {
					modal := a.DeleteVaultModal(actualIdx)
					a.App.SetRoot(modal, false)
				}
			}
		case 'u':
			currentIndex := a.VaultList.GetCurrentItem()
			if currentIndex >= 0 && currentIndex < len(a.Vault) {
				entry := a.FilteredVault[currentIndex]
				actualIdx := a.findFilteredVaultIndex(entry)

				if actualIdx >= 0 {
					flex := a.UpdateVaultModal(actualIdx)
					a.App.SetRoot(flex, true)
				}
			}
		case 'g':
			generatedPassword, err := vault.GeneratePassword(20, vault.DefaultChars)
			if err != nil {
				modal := a.ErrorModal(err.Error(), a.Root)
				a.App.SetRoot(modal, true)
			}
			modal := a.GeneratedPasswordModal(string(generatedPassword))
			a.App.SetRoot(modal, true)
		case 'b':
			backupModal := a.BackupModal()
			a.App.SetRoot(backupModal, true)
		case 'l':
			if a.ToggleShowBackup {
				backups, err := a.ListBackupsFlex()
				if err != nil {
					modal := a.ErrorModal(err.Error(), a.Root)
					a.App.SetRoot(modal, true)
					return nil
				}
				a.App.SetRoot(backups, true)
			} else {
				a.App.SetRoot(a.Root, true)
			}
		case 'c':
			currentIndex := a.VaultList.GetCurrentItem()
			if currentIndex >= 0 && currentIndex < len(a.Vault) {
				entry := a.FilteredVault[currentIndex]
				actualIdx := a.findFilteredVaultIndex(entry)

				if actualIdx >= 0 {
					successModal := a.CopyDirectlyToClipboard(actualIdx)
					a.App.SetRoot(successModal, true)
				}
			}

		case 'q':
			a.App.Stop()
		case '\t':
			a.App.SetFocus(a.SearchInput)
			return nil
		}

		return event
	})

	a.VaultList.SetSelectedFunc(func(itemIdx int, primaryText, secondaryText string, _ rune) {
		if itemIdx >= 0 && itemIdx < len(a.FilteredVault) {
			modal := a.ModalVaultInfoByVault(a.FilteredVault[itemIdx])
			a.App.SetRoot(modal, false)
		}
	})
}

// SyncFilteredVault ensures that if there is text in the SearchInput, the
// FilteredVault is updated with what the entries should be
func (a *App) SyncFilteredVault() {
	text := a.SearchInput.GetText()
	a.FilteredVault = FilterVaultEntries(a.Vault, text)
}

// FilterVaultEntries contains the logic of filtering the vault list with each
// keypress while searching
func FilterVaultEntries(vault []model.VaultEntry, searchText string) []model.VaultEntry {
	if searchText == "" {
		return vault
	}
	filtered := make([]model.VaultEntry, 0)

	for _, v := range vault {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(searchText)) {
			entry := v
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

// findFilteredVaultIndex finds the actual vault entry after search
func (a *App) findFilteredVaultIndex(entry model.VaultEntry) int {
	for i, v := range a.Vault {
		if v.UpdatedAt == entry.UpdatedAt {
			return i
		}
	}
	return -1
}

// RefreshRoot refreshes the application root.
func (a *App) RefreshRoot() {
	help := tview.NewTextView().
		SetText(HelpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.SearchBar, 3, 1, false).
		AddItem(a.VaultListView(), 0, 1, true).
		AddItem(help, 3, 1, false)
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.App.Stop()
			return nil
		}
		return event
	})

	a.Root = root
}

// SaveVault encrypts the vault and saves the vault to disk
func (a *App) SaveVault() {
	defer func() {
		a.PopulateVaultList()
	}()
	encryptedCipherText, err := crypt.EncryptVault(a.Vault, a.Keyring)
	if err != nil {
		modal := a.ErrorModal(fmt.Sprintf("Failed to save vault: %v", err), a.Root)
		a.App.SetRoot(modal, true)
	}

	if err := utils.WriteToFile(a.VaultFile.Name(), model.FileVault, encryptedCipherText); err != nil {
		modal := a.ErrorModal(fmt.Sprintf("Failed to save vault: %v", err), a.Root)
		a.App.SetRoot(modal, true)
		return
	}

	// need to refresh the app.VaultFile
	a.VaultFile, err = os.OpenFile(a.VaultFile.Name(), os.O_RDWR, 0o600)
	if err != nil {
		modal := a.ErrorModal(fmt.Sprintf("Failed to open vault file: %v", err), a.Root)
		a.App.SetRoot(modal, true)
		return
	}
}

// VaultListView builds the list view of the VaultList in a Flex primitive
func (a *App) VaultListView() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(a.VaultList, 0, 1, true).
		AddItem(box, 0, 1, false)
}
