package cmd

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var helpText = "Commands"

// inputPassword string
// passwordInput *tview.InputField
//
// home, _     = os.UserHomeDir()
// CONFIG_PATH = path.Join(home, ".config", "gopass")
// CONFIG_FILE = path.Join(CONFIG_PATH, "gopass-cfg.json")

type App struct {
	App   *tview.Application
	Vault []model.VaultEntry
	Cfg   model.Config

	VaultList *tview.List
	Root      *tview.Flex
}

func (a *App) PopulateVaultList() {
	a.VaultList = tview.NewList()
	for _, v := range a.Vault {
		a.VaultList.AddItem(v.Name, "", 0, nil)
	}

	a.VaultList.SetBorder(true)
	a.VaultList.SetTitle("Vault")
	a.VaultList.SetBackgroundColor(tcell.ColorBlack)
	a.VaultList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// TODO: Add keys to proceed with vault actions
		switch event.Rune() {
		// case 'a':
		// 	AddToVault()
		}
		return event
	})

	a.VaultList.SetSelectedFunc(func(itemIdx int, primaryText, secondaryText string, _ rune) {
		modal := a.ModalVaultInfo(itemIdx)
		a.App.SetRoot(modal, false)
	})
}

func (a *App) VaultListView() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(a.VaultList, 0, 1, true).
		AddItem(box, 0, 1, false)
}

func (a *App) ModalVaultInfo(idx int) *tview.Modal {
	text := fmt.Sprintf(`
	Name: %s
	Password: %s
	Notes: %s
	`, a.Vault[idx].Name, crypt.DecryptPassword(a.Vault[idx].Password), a.Vault[idx].Notes)
	modal := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle("Vault Info")
	modal.SetText(text)
	modal.SetBorder(true)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		a.App.SetRoot(a.Root, true)
	})

	return modal
}

func NewApp() *App {
	return &App{}
}

func TviewRun() {
	app := NewApp()
	app.App = tview.NewApplication()
	cfg, err := utils.CheckConfig("")
	if err != nil {
		panic(err)
	}
	app.Cfg = cfg

	vaultF, _ := utils.OpenVault(cfg.VaultName)
	vault := crypt.DecryptVault(vaultF)
	app.Vault = vault

	app.PopulateVaultList()

	help := tview.NewTextView().
		SetText(helpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.VaultListView(), 0, 1, true).
		AddItem(help, 3, 1, false)
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.App.Stop()
			return nil
		}
		return event
	})
	app.Root = root

	if err := app.App.SetRoot(app.Root, true).Run(); err != nil {
		panic(err)
	}
}

// func isUserLoggedIn(cfg model.Config) bool {
// 	now := time.Now().UnixMilli()
// 	return !utils.IsAccessBeforeLogin(cfg, now)
// }
