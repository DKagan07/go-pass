package components_tview

import "github.com/rivo/tview"

func PasswordInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Password").
		SetMaskCharacter('*')
}
