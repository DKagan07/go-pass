package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/crypt"
	"go-pass/utils"
)

var helpText = "Commands"

// inputPassword string
// passwordInput *tview.InputField
//
// home, _     = os.UserHomeDir()
// CONFIG_PATH = path.Join(home, ".config", "gopass")
// CONFIG_FILE = path.Join(CONFIG_PATH, "gopass-cfg.json")

func TviewRun() {
	app := tview.NewApplication()
	cfg, err := utils.CheckConfig("")
	if err != nil {
		panic(err)
	}

	// Problem here v
	vaultF, _ := utils.OpenVault(cfg.VaultName)
	vault := crypt.DecryptVault(vaultF)

	l := tview.NewList()
	for _, v := range vault {
		l.AddItem(v.Name, "", 0, nil)
	}

	l.SetBorder(true)
	l.SetTitle("Vault")
	l.SetBackgroundColor(tcell.ColorBlack)

	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	middleL := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(l, 0, 1, true).
		AddItem(box, 0, 1, false)

	help := tview.NewTextView().
		SetText(helpText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	container := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(middleL, 0, 1, true).
		AddItem(help, 3, 1, false)
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(container, true).Run(); err != nil {
		panic(err)
	}
}

// func isUserLoggedIn(cfg model.Config) bool {
// 	now := time.Now().UnixMilli()
// 	return !utils.IsAccessBeforeLogin(cfg, now)
// }
