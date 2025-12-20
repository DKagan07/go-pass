package tui

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
)

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
		vError := ValidateAddInput(formName, formUsername, formPassword)
		if vError != nil {
			eModal := a.ErrorModal(vError.Error(), a.Root)
			a.App.SetRoot(eModal, true)
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

func ValidateAddInput(name, username, password string) error {
	if strings.EqualFold(name, "") {
		return &ValidationError{Field: "Name", Message: "Name cannot be empty"}
	}

	if strings.EqualFold(username, "") {
		return &ValidationError{Field: "Username", Message: "Username cannot be empty"}
	}

	if strings.EqualFold(password, "") {
		return &ValidationError{Field: "Password", Message: "Password cannot be empty"}
	}

	return nil
}
