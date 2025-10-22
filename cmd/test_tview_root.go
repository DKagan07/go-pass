package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var helpText = " a: Add | d: Delete | u: Update | tab: Switch between Search and Vault "

type App struct {
	App       *tview.Application
	VaultFile *os.File
	Vault     []model.VaultEntry
	Cfg       model.Config

	VaultList *tview.List
	Root      *tview.Flex
	SearchBar *tview.Flex
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
	a.VaultList.SetTitle(" Vault ")
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
		case 'u':
			currentIndex := a.VaultList.GetCurrentItem()
			if currentIndex >= 0 && currentIndex < len(a.Vault) {
				flex := a.UpdateVaultModal(currentIndex)
				a.App.SetRoot(flex, true)
			}
		case '\t':
			a.App.SetFocus(a.SearchBar)
		}

		return event
	})

	a.VaultList.SetSelectedFunc(func(itemIdx int, primaryText, secondaryText string, _ rune) {
		modal := a.ModalVaultInfo(itemIdx)
		a.App.SetRoot(modal, false)
	})
}

func (a *App) CreateSearchBar() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	search := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldBackgroundColor(tcell.ColorBlack)
	search.SetBackgroundColor(tcell.ColorBlack)
	search.SetChangedFunc(func(text string) {
		a.VaultList.Clear()
		for _, v := range a.Vault {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(text)) {
				a.VaultList.AddItem(v.Name, "", 0, nil)
			}
		}
	})

	search.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			itemIdx := a.VaultList.GetCurrentItem()
			modal := a.ModalVaultInfo(itemIdx)
			a.App.SetRoot(modal, false)
		}
	})

	search.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			a.App.SetFocus(a.VaultListView())
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

	modal.SetTitle(" Vault Info ")
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

		// Validation
		if strings.EqualFold(formName, "") {
			modal := a.ErrorModal("Name cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		if strings.EqualFold(formUsername, "") {
			modal := a.ErrorModal("Username cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		if strings.EqualFold(formPassword, "") {
			modal := a.ErrorModal("Password cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		a.AddToVault(formName, formNotes, formUsername, formPassword)

		a.PopulateVaultList()
		a.RefreshRoot()
		a.App.SetRoot(a.Root, true)
	})
	inputForm.SetTitle(" Add Vault ")
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

	modal.SetTitle(" Delete Vault ")
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}

func (a *App) DeleteFromVault(vaultIdx int) {
	a.Vault = slices.Delete(a.Vault, vaultIdx, vaultIdx+1)
	a.SaveVault()
}

func (a *App) UpdateVaultModal(currIdx int) *tview.Flex {
	entry := a.Vault[currIdx]

	form := tview.NewForm().
		AddInputField("Name", entry.Name, 0, nil, nil).
		AddInputField("Username", entry.Username, 0, nil, nil).
		AddInputField("Password", string(crypt.DecryptPassword(entry.Password)), 0, nil, nil).
		AddInputField("Notes", entry.Notes, 0, nil, nil)
	form.AddButton("Save", func() {
		formName := form.GetFormItem(0).(*tview.InputField).GetText()
		formUsername := form.GetFormItem(1).(*tview.InputField).GetText()
		formPassword := form.GetFormItem(2).(*tview.InputField).GetText()
		formNotes := form.GetFormItem(3).(*tview.InputField).GetText()

		if strings.EqualFold(formName, "") {
			modal := a.ErrorModal("Name cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		if strings.EqualFold(formUsername, "") {
			modal := a.ErrorModal("Username cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		if strings.EqualFold(formPassword, "") {
			modal := a.ErrorModal("Password cannot be empty", a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		newEntry := model.VaultEntry{
			Name:      formName,
			Username:  formUsername,
			Notes:     formNotes,
			Password:  crypt.EncryptPassword([]byte(formPassword)),
			UpdatedAt: entry.UpdatedAt,
		}

		a.UpdateVaultEntry(currIdx, newEntry)
		a.PopulateVaultList()
		a.RefreshRoot()
		a.App.SetRoot(a.Root, true)
	})

	form.SetTitle(" Update Vault ")
	form.SetBorder(true)
	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorBlack)
	form.SetButtonBackgroundColor(tcell.Color103)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(nil, 0, 1, false)

	return flex
}

func (a *App) UpdateVaultEntry(currIdx int, newEntry model.VaultEntry) {
	a.Vault[currIdx] = newEntry
	a.SaveVault()
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

func (a *App) ErrorModal(errMsg string, dest tview.Primitive) *tview.Modal {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"OK"}).
		SetButtonBackgroundColor(tcell.ColorBlack).
		SetText(errMsg).
		SetTextColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.App.SetRoot(dest, true)
		})

	modal.SetTitle(" Error! ")
	modal.SetTitleColor(tcell.ColorRed)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	return modal
}

func NewApp() *App {
	return &App{}
}

func TviewRun() {
	app := NewApp()
	app.App = tview.NewApplication()
	cfgFile, ok, err := utils.OpenConfig("")
	if ok && err == nil {
		panic(errors.New("a file is not found. need to 'init'"))
		// TODO: implement 'init'
	}
	cfg := crypt.DecryptConfig(cfgFile)
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
		AddItem(app.CreateSearchBar(), 3, 1, true).
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

	now := time.Now().UnixMilli()
	if !utils.IsAccessBeforeLogin(cfg, now) {
		loginForm := tview.NewForm().
			AddPasswordField("Master Password", "", 0, '*', nil)

		loginForm.SetTitle(" Login ")
		loginForm.SetBorder(true)
		loginForm.SetBackgroundColor(tcell.ColorBlack)
		loginForm.SetFieldBackgroundColor(tcell.ColorBlack)
		loginForm.SetButtonBackgroundColor(tcell.Color103)
		loginForm.AddButton("Login", func() {
			masterPassword := loginForm.GetFormItem(0).(*tview.InputField).GetText()
			// if err is not nil, then the user has input the wrong password
			err := bcrypt.CompareHashAndPassword(cfg.MasterPassword, []byte(masterPassword))
			if err != nil {
				modal := app.ErrorModal("Incorrect Master Password", loginForm)
				loginForm.GetFormItem(0).(*tview.InputField).SetText("")
				app.App.SetRoot(modal, false)
				return
			}

			cfg.LastVisited = now
			encryptedCfg, err := crypt.EncryptConfig(cfg)
			if err != nil {
				panic(err)
			}
			utils.WriteToFile(cfgFile, encryptedCfg)
			app.App.SetRoot(app.Root, true)
		})

		loginPage := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(loginForm, 0, 1, true).
			AddItem(help, 3, 1, false)

		if err := app.App.SetRoot(loginPage, true).Run(); err != nil {
			panic(err)
		}
	} else {
		if err := app.App.SetRoot(app.Root, true).Run(); err != nil {
			panic(err)
		}
	}
}
