package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var helpText = "a: Add | Commands"

// inputPassword string
// passwordInput *tview.InputField
//
// home, _     = os.UserHomeDir()
// CONFIG_PATH = path.Join(home, ".config", "gopass")
// CONFIG_FILE = path.Join(CONFIG_PATH, "gopass-cfg.json")

type App struct {
	App       *tview.Application
	VaultFile *os.File
	Vault     []model.VaultEntry
	Cfg       model.Config

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
		case 'a':
			a.ModalAddVault()
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

// TODO: Figure out which is better, returning a tview.Primitive, or nothing and
// handling setting the root in the caller.
// TODO: Figure out how to do this in a modal? Or maybe flex is just better.
// And get this fucntionality working

func (a *App) ModalAddVault() {
	inputForm := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddInputField("Notes", "", 0, nil, nil)

	inputForm.AddButton("Save", func() {
		formName := inputForm.GetFormItemByLabel("Name").GetLabel()
		formUsername := inputForm.GetFormItemByLabel("Username").GetLabel()
		formPassword := inputForm.GetFormItemByLabel("Password").GetLabel()
		formNotes := inputForm.GetFormItemByLabel("Notes").GetLabel()

		fmt.Println("formName: ", formName)
		fmt.Println("formUsername: ", formUsername)
		fmt.Println("formPassword: ", formPassword)
		fmt.Println("formNotes: ", formNotes)

		a.App.SetRoot(a.Root, true)
	})
	inputForm.SetTitle("Add Vault")
	inputForm.SetBorder(true)
	inputForm.SetBackgroundColor(tcell.ColorBlack)
	inputForm.SetFieldBackgroundColor(tcell.ColorBlack)
	inputForm.SetButtonBackgroundColor(tcell.Color103)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(inputForm, 0, 1, true).
		AddItem(nil, 0, 1, false)

	a.App.SetRoot(flex, true)
}

func (a *App) AddToVault(name, notes, username string, password []byte) {
	now := time.Now().UnixMilli()
	vault := model.VaultEntry{
		Name:      name,
		Username:  username,
		Notes:     notes,
		Password:  password,
		UpdatedAt: now,
	}

	a.Vault = append(a.Vault, vault)

	encryptedCipherText, err := crypt.EncryptVault(a.Vault)
	if err != nil {
		panic(err)
	}
	utils.WriteToFile(a.VaultFile, encryptedCipherText)
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
	app.VaultFile = vaultF
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
