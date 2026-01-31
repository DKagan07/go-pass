package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/model"
)

var HelpText = "a: Add | d: Delete | u: Update | g: Generate Password | c: Copy | b: Backup | l: Toggle Backup Display | q: Quit | tab: Switch Search and Vault"

// App is the structure that controls all the actions for the TUI
type App struct {
	App              *tview.Application
	VaultFile        *os.File
	Vault            []model.VaultEntry
	FilteredVault    []model.VaultEntry
	Cfg              *model.Config
	Keyring          *model.MasterAESKeyManager
	NumRetries       int32
	ToggleShowBackup bool

	VaultList   *tview.List
	Root        *tview.Flex
	SearchBar   *tview.Flex
	SearchInput *tview.InputField
}

// NewApp returns a new App
func NewApp() *App {
	return &App{
		ToggleShowBackup: true,
	}
}

// CreateRoot returns a Flex primitive of the root of the program.
// The root is a combination of the search bar and the list model.VaultEntry
func (a *App) CreateRoot() *tview.Flex {
	help := tview.NewTextView().
		SetText(HelpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.CreateSearchBar(), 3, 1, true).
		AddItem(a.VaultListView(), 0, 1, true).
		AddItem(help, 3, 1, false)
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.App.Stop()
			return nil
		}
		return event
	})

	return root
}

// TviewRun is the main entry point to the TUI part of the program. It sets up
// the app and login screen.
func TviewRun() {
	app := NewApp()
	app.App = tview.NewApplication()
	loginPage := app.LoginPage()

	if err := app.App.SetRoot(loginPage, true).Run(); err != nil {
		modal := app.ExitErrorModal(err.Error())
		app.App.SetRoot(modal, true)
	}
}
