package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/model"
	"go-pass/utils"
)

// LoginPage holds the functionality to create the login page, when the user
// first runs the app with the `-o` flag
func (a *App) LoginPage() *tview.Flex {
	var loginPage *tview.Flex

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
		a.Keyring = keyring

		cfgFile, ok, err := utils.OpenConfig("")
		if ok && err == nil {
			modal := a.ExitErrorModal("a file is not found. need to run 'gopass init'")
			a.App.SetRoot(modal, true)
		}

		cfg, err := crypt.DecryptConfig(cfgFile, a.Keyring, false)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "decrypting contents") {
				a.NumRetries++
				errMsg = fmt.Sprintf("Authentication Failed, retries: %d", a.NumRetries)
			}

			if a.NumRetries >= 3 {
				finalModal := a.ExitErrorModal("All attempts made")
				a.App.SetRoot(finalModal, true)
				return
			} else {
				modal := a.ErrorModal(errMsg, loginPage)
				a.App.SetRoot(modal, true)
				return
			}
		}

		a.Cfg = cfg

		vaultF, err := utils.OpenVault(cfg.VaultName)
		if err != nil {
			modal := a.ErrorModal(err.Error(), loginPage)
			a.App.SetRoot(modal, true)
			return
		}

		a.VaultFile = vaultF
		vault, err := crypt.DecryptVault(vaultF, a.Keyring, false)
		if err != nil {
			modal := a.ExitErrorModal(err.Error())
			a.App.SetRoot(modal, true)
		}

		a.Vault = vault
		a.FilteredVault = vault

		a.PopulateVaultList()

		root := a.CreateRoot()
		a.Root = root

		cfg.LastVisited = time.Now().UnixMilli()
		encryptedCfg, err := crypt.EncryptConfig(cfg, a.Keyring)
		if err != nil {
			modal := a.ExitErrorModal(err.Error())
			a.App.SetRoot(modal, true)
		}
		utils.WriteToFile(cfgFile, encryptedCfg)

		a.App.SetRoot(a.Root, true)
	})

	loginPage = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 1, true)

	return loginPage
}
