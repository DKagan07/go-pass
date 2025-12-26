package tui

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
)

// UpdateVaultModal returns a Flex primitive containing the modal and form for
// a user to update their vault entry
func (a *App) UpdateVaultModal(currIdx int) *tview.Flex {
	entry := a.Vault[currIdx]

	decryptedPass, err := crypt.DecryptPassword(entry.Password, a.Keyring, false)
	if err != nil {
		modal := a.ErrorModal(err.Error(), a.Root)
		a.App.SetRoot(modal, true)
	}

	form := tview.NewForm().
		AddInputField("Name", entry.Name, 0, nil, nil).
		AddInputField("Username", entry.Username, 0, nil, nil).
		AddInputField("Password", decryptedPass, 0, nil, nil).
		AddInputField("Notes", entry.Notes, 0, nil, nil)
	form.AddButton("Save", func() {
		formName := form.GetFormItem(0).(*tview.InputField).GetText()
		formUsername := form.GetFormItem(1).(*tview.InputField).GetText()
		formPassword := form.GetFormItem(2).(*tview.InputField).GetText()
		formNotes := form.GetFormItem(3).(*tview.InputField).GetText()

		newEntry, err := a.ValidateUpdateInputs(
			formName,
			formUsername,
			formPassword,
			formNotes,
		)
		if err != nil {
			modal := a.ErrorModal(err.Error(), a.Root)
			a.App.SetRoot(modal, false)
			return
		}

		a.UpdateVaultEntry(currIdx, *newEntry)
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

// UpdateVaultEntry contains the business logic of updating the vault on disk
func (a *App) UpdateVaultEntry(currIdx int, newEntry model.VaultEntry) {
	a.Vault[currIdx] = newEntry
	a.SaveVault()
}

// ValidateUpdateInputs ensures that the necessary inputs are present when
// updating a vault entry
func (a *App) ValidateUpdateInputs(
	name, username, password, notes string,
) (*model.VaultEntry, error) {
	now := time.Now().UnixMilli()

	if strings.EqualFold(name, "") {
		return nil, &ValidationError{Field: "Name", Message: "Name cannot be empty"}
	}

	if strings.EqualFold(username, "") {
		return nil, &ValidationError{Field: "Username", Message: "Username cannot be empty"}
	}

	if strings.EqualFold(password, "") {
		return nil, &ValidationError{Field: "Password", Message: "Password cannot be empty"}
	}

	p, err := crypt.EncryptPassword([]byte(password), a.Keyring)
	if err != nil {
		return nil, &ValidationError{
			Field:   "EncryptedPassword",
			Message: "Failed to encrypt password",
		}
	}

	return &model.VaultEntry{
		Name:      name,
		Username:  username,
		Notes:     notes,
		Password:  []byte(p),
		UpdatedAt: now,
	}, nil
}
