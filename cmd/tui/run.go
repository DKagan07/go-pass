package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var HelpText = "a: Add | d: Delete | u: Update | g: Generate Password | b: Create backup | l: Toggle Backup Display | q: Quit | tab: Switch Search and Vault"

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
	var loginPage *tview.Flex

	// if !app.IsLoggedIn {
	loginForm := tview.NewForm().
		AddPasswordField("Master Password", "", 0, '*', nil)

	loginForm.SetTitle(" Login ")
	loginForm.SetBorder(true)
	loginForm.SetBackgroundColor(tcell.ColorBlack)
	loginForm.SetFieldBackgroundColor(tcell.ColorBlack)
	loginForm.SetButtonBackgroundColor(tcell.Color103)
	loginForm.AddButton("Login", func() {
		masterPassword := loginForm.GetFormItem(0).(*tview.InputField).GetText()

		keyring := model.NewMasterAESKeyManager(masterPassword)
		app.Keyring = keyring

		cfgFile, ok, err := utils.OpenConfig("")
		if ok && err == nil {
			modal := app.ExitErrorModal("a file is not found. need to run 'gopass init'")
			app.App.SetRoot(modal, true)
		}

		cfg, err := crypt.DecryptConfig(cfgFile, app.Keyring, false)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "decrypting contents") {
				app.NumRetries++
				errMsg = fmt.Sprintf("Authentication Failed, retries: %d", app.NumRetries)
			}

			if app.NumRetries >= 3 {
				finalModal := app.ExitErrorModal("All attempts made")
				app.App.SetRoot(finalModal, true)
				return
			} else {
				modal := app.ErrorModal(errMsg, loginPage)
				app.App.SetRoot(modal, true)
				return
			}
		}

		app.Cfg = cfg

		vaultF, err := utils.OpenVault(cfg.VaultName)
		if err != nil {
			modal := app.ErrorModal(err.Error(), loginPage)
			app.App.SetRoot(modal, true)
			return
		}

		app.VaultFile = vaultF
		vault, err := crypt.DecryptVault(vaultF, app.Keyring, false)
		if err != nil {
			modal := app.ExitErrorModal(err.Error())
			app.App.SetRoot(modal, true)
		}

		app.Vault = vault
		app.FilteredVault = vault

		app.PopulateVaultList()

		root := app.CreateRoot()
		app.Root = root

		cfg.LastVisited = time.Now().UnixMilli()
		encryptedCfg, err := crypt.EncryptConfig(cfg, app.Keyring)
		if err != nil {
			modal := app.ExitErrorModal(err.Error())
			app.App.SetRoot(modal, true)
		}
		utils.WriteToFile(cfgFile, encryptedCfg)

		app.App.SetRoot(app.Root, true)
	})

	loginPage = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 1, true)

	if err := app.App.SetRoot(loginPage, true).Run(); err != nil {
		modal := app.ExitErrorModal(err.Error())
		app.App.SetRoot(modal, true)
	}
}
