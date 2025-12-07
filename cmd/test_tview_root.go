package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/cmd/vault"
	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

var helpText = " a: Add | d: Delete | u: Update | g: Generate Password | tab: Switch between Search and Vault "

type App struct {
	App           *tview.Application
	VaultFile     *os.File
	Vault         []model.VaultEntry
	FilteredVault []model.VaultEntry
	Cfg           model.Config
	Keyring       *model.MasterAESKeyManager

	VaultList   *tview.List
	Root        *tview.Flex
	SearchBar   *tview.Flex
	SearchInput *tview.InputField
}

func (a *App) PopulateVaultList() {
	// Alphebetize the vault by name
	sort.Slice(a.Vault, func(i, j int) bool {
		return a.Vault[i].Name < a.Vault[j].Name
	})

	a.VaultList = tview.NewList()
	for _, v := range a.Vault {
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
		// TODO: Add keys to proceed with vault actions
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
			generatedPassword := vault.GeneratePassword(20, vault.DefaultChars)
			modal := a.GeneratedPasswordModal(string(generatedPassword))
			a.App.SetRoot(modal, true)
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

// findFilteredVaultIndex finds the actual vault entry after search
func (a *App) findFilteredVaultIndex(entry model.VaultEntry) int {
	for i, v := range a.Vault {
		if v.UpdatedAt == entry.UpdatedAt {
			return i
		}
	}
	return -1
}

func (a *App) ModalVaultInfoByVault(ve model.VaultEntry) *tview.Modal {
	decryptedPassword := crypt.DecryptPassword(ve.Password, a.Keyring, false)
	text := fmt.Sprintf(`
	Name: %s
	Username: %s
	Password: %s
	Notes: %s
	`, ve.Name, ve.Username, decryptedPassword, ve.Notes)
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

		for _, v := range a.Vault {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(text)) {
				vault := v
				a.FilteredVault = append(a.FilteredVault, vault)

				a.VaultList.AddItem(vault.Name, "", 0, func() {
					m := a.ModalVaultInfoByVault(vault)
					a.App.SetRoot(m, false)
				})
			}
		}
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

func (a *App) VaultListView() *tview.Flex {
	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(a.VaultList, 0, 1, true).
		AddItem(box, 0, 1, false)
}

func (a *App) ModalVaultInfoByIdx(idx int) *tview.Modal {
	entry := a.Vault[idx]
	decryptedPassword := crypt.DecryptPassword(entry.Password, a.Keyring, false)
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

	inputForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			a.App.SetRoot(a.Root, true)
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(inputForm, 0, 1, true).
		AddItem(nil, 0, 1, false)

	return flex
}

func (a *App) AddToVault(name, notes, username, password string) {
	passwordBytes := []byte(password)
	encryptedPassword, _ := crypt.EncryptPassword(passwordBytes, a.Keyring)
	now := time.Now().UnixMilli()

	vault := model.VaultEntry{
		Name:      name,
		Username:  username,
		Notes:     notes,
		Password:  []byte(encryptedPassword),
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
		AddInputField("Password", crypt.DecryptPassword(entry.Password, a.Keyring, false), 0, nil, nil).
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

		p, _ := crypt.EncryptPassword([]byte(formPassword), a.Keyring)
		newEntry := model.VaultEntry{
			Name:      formName,
			Username:  formUsername,
			Notes:     formNotes,
			Password:  []byte(p),
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

func (a *App) GeneratedPasswordModal(generatedPass string) *tview.Modal {
	modal := tview.NewModal().
		AddButtons([]string{"OK", "Copy"}).
		SetBackgroundColor(tcell.ColorBlack)

	modal.SetTitle(" Generated Password ")
	modal.SetText(generatedPass)
	modal.SetBorder(true)
	modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if strings.EqualFold(buttonLabel, "Copy") {
			err := clipboard.WriteAll(generatedPass)
			if err != nil {
				panic(err)
			}
		}
		a.App.SetRoot(a.Root, true)
	})

	return modal
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
		panic(err)
	}

	utils.WriteToFile(a.VaultFile, encryptedCipherText)
}

func (a *App) CreateRoot() *tview.Flex {
	help := tview.NewTextView().
		SetText(helpText).
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
			panic(errors.New("a file is not found. need to 'init'"))
			// TODO: implement 'init'
		}
		cfg := crypt.DecryptConfig(cfgFile, app.Keyring, false)
		app.Cfg = cfg

		vaultF, _ := utils.OpenVault(cfg.VaultName)
		app.VaultFile = vaultF
		vault := crypt.DecryptVault(vaultF, app.Keyring, false)
		app.Vault = vault
		app.FilteredVault = vault

		app.PopulateVaultList()

		root := app.CreateRoot()
		app.Root = root

		cfg.LastVisited = time.Now().UnixMilli()
		encryptedCfg, err := crypt.EncryptConfig(cfg, app.Keyring)
		if err != nil {
			panic(err)
		}
		utils.WriteToFile(cfgFile, encryptedCfg)

		app.App.SetRoot(app.Root, true)
	})

	loginPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 1, true)

	if err := app.App.SetRoot(loginPage, true).Run(); err != nil {
		panic(err)
	}
}
