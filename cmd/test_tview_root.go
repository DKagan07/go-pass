package cmd

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var helpText = "a: Add | d: Delete"

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
	// Alphebetize the vault by name
	sort.Slice(a.Vault, func(i, j int) bool {
		return a.Vault[i].Name < a.Vault[j].Name
	})

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
			flex := a.ModalAddVault()
			a.App.SetRoot(flex, true)
		case 'd':
			currentIndex := a.VaultList.GetCurrentItem()
			if currentIndex >= 0 && currentIndex < len(a.Vault) {
				modal := a.DeleteVaultModal(currentIndex)
				a.App.SetRoot(modal, false)
			}
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
	entry := a.Vault[idx]
	decryptedPassword := crypt.DecryptPassword(entry.Password)
	text := fmt.Sprintf(`
	Name: %s
	Username: %s
	Password: %s
	Notes: %s
	`, entry.Name, entry.Username, decryptedPassword, entry.Notes)
	modal := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle("Vault Info")
	modal.SetText(text)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		a.App.SetRoot(a.Root, true)
	})

	return modal
}

func (a *App) ModalAddVault() *tview.Flex {
	inputForm := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddInputField("Notes", "", 0, nil, nil)

	inputForm.AddButton("Save", func() {
		formName := inputForm.GetFormItem(0).(*tview.InputField).GetText()
		formUsername := inputForm.GetFormItem(1).(*tview.InputField).GetText()
		formPassword := inputForm.GetFormItem(2).(*tview.InputField).GetText()
		formNotes := inputForm.GetFormItem(3).(*tview.InputField).GetText()

		// TODO: Really need to add some validation to make sure that:
		// 1. Name is not empty and unique
		// 2. Username is not empty
		// 3. Password is not empty

		a.AddToVault(formName, formNotes, formUsername, formPassword)

		a.PopulateVaultList()
		a.RefreshRoot()
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

	return flex
}

func (a *App) AddToVault(name, notes, username, password string) {
	passwordBytes := []byte(password)
	encryptedPassword := crypt.EncryptPassword(passwordBytes)
	now := time.Now().UnixMilli()

	vault := model.VaultEntry{
		Name:      name,
		Username:  username,
		Notes:     notes,
		Password:  encryptedPassword,
		UpdatedAt: now,
	}
	a.Vault = append(a.Vault, vault)
	a.SaveVault()
}

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

	modal.SetTitle("Delete Vault")
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}

func (a *App) DeleteFromVault(vaultIdx int) {
	a.Vault = slices.Delete(a.Vault, vaultIdx, vaultIdx+1)
	a.SaveVault()
	a.PopulateVaultList()
}

func (a *App) RefreshRoot() {
	help := tview.NewTextView().
		SetText(helpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
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
