package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *App) CreateSearchBar() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	search := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldBackgroundColor(tcell.ColorBlack)
	search.SetBackgroundColor(tcell.ColorBlack)
	a.SearchInput = search

	search.SetChangedFunc(func(text string) {
		a.VaultList.Clear()
		a.FilteredVault = nil

		a.FilterVaultBySearch(text)
	})

	search.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			a.App.SetFocus(a.VaultList)
			return nil
		}
		return event
	})

	search.SetBorder(true)

	searchbar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(search, 0, 1, true).
		AddItem(box, 0, 1, false)

	a.SearchBar = searchbar
	return searchbar
}

func (a *App) FilterVaultBySearch(searchText string) {
	a.VaultList.Clear()
	a.FilteredVault = nil

	for _, v := range a.Vault {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(searchText)) {
			vault := v
			a.FilteredVault = append(a.FilteredVault, vault)

			a.VaultList.AddItem(vault.Name, "", 0, func() {
				m := a.ModalVaultInfoByVault(vault)
				a.App.SetRoot(m, false)
			})
		}
	}
}
