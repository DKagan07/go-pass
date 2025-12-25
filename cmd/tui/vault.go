package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/cmd/vault"
	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

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

func (a *App) SyncFilteredVault() {
	text := a.SearchInput.GetText()
	a.FilteredVault = FilterVaultEntries(a.Vault, text)
}

func FilterVaultEntries(vault []model.VaultEntry, searchText string) []model.VaultEntry {
	if searchText == "" {
		return vault
	}
	filtered := make([]model.VaultEntry, 0)

	for _, v := range vault {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(searchText)) {
			vault := v
			filtered = append(filtered, vault)
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

func (a *App) RefreshRoot() {
	help := tview.NewTextView().
		SetText(helpText).
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

func (a *App) SaveVault() {
	encryptedCipherText, err := crypt.EncryptVault(a.Vault, a.Keyring)
	if err != nil {
		modal := a.ErrorModal(fmt.Sprintf("Failed to save vault: %v", err), a.Root)
		a.App.SetRoot(modal, true)
	}

	utils.WriteToFile(a.VaultFile, encryptedCipherText)
}

func (a *App) VaultListView() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(a.VaultList, 0, 1, true).
		AddItem(box, 0, 1, false)
}
