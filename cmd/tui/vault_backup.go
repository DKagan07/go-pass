package tui

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"go-pass/cmd/vault"
	"go-pass/utils"
)

func (a *App) BackupModal() *tview.Modal {
	backupModal := tview.NewModal().
		AddButtons([]string{"Yes", "No"}).
		SetBackgroundColor(tcell.ColorBlack)

	backupModal.SetTitle(" Confirm Backup ")
	backupModal.SetBorder(true)
	backupModal.SetText("Do you want to create a current backup?")
	backupModal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	backupModal.SetDoneFunc(func(buttonIdx int, buttonLabel string) {
		if strings.EqualFold(buttonLabel, "Yes") {
			now := time.Now()
			successMsg, err := vault.BackupVault("", a.Cfg.VaultName, "", now, a.Keyring)
			if err != nil {
				errModal := a.ErrorModal(err.Error(), a.Root)
				a.App.SetRoot(errModal, true)
			}

			successModal := tview.NewModal().
				AddButtons([]string{"Ok"}).
				SetBackgroundColor(tcell.ColorBlack)
			successModal.SetBorder(true)
			successModal.SetText(successMsg)
			successModal.SetTextColor(tcell.ColorGreen)
			successModal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
			successModal.SetDoneFunc(func(_ int, buttonLabel string) {
				if strings.EqualFold(buttonLabel, "Ok") {
					a.App.SetRoot(a.Root, true)
				}
			})

			a.RefreshRoot()
			a.App.SetRoot(successModal, true)
		} else if strings.EqualFold(buttonLabel, "No") {
			a.App.SetRoot(a.Root, true)
		}
	})

	return backupModal
}

func (a *App) ListBackupsFlex() (*tview.Flex, error) {
	a.ToggleShowBackup = !a.ToggleShowBackup

	dirEntries, err := os.ReadDir(utils.BACKUP_DIR)
	if err != nil {
		return nil, err
	}

	if len(dirEntries) == 0 {
		modal := tview.NewModal().
			AddButtons([]string{"Ok"}).
			SetBackgroundColor(tcell.ColorBlack)

		modal.SetBorderStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
		modal.SetText("No backups found!")
		modal.SetDoneFunc(func(bIdx int, _ string) {
			a.App.SetRoot(a.Root, true)
		})
		return tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(modal, 0, 1, true), nil
	}

	help := tview.NewTextView().
		SetText(" l: Toggle back to Vault ").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	help.SetBorder(true).SetTitle(" Help ")

	backupRoot := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.BackupListView(dirEntries), 0, 1, true).
		AddItem(help, 3, 1, false)
	backupRoot.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'l':
			a.ToggleShowBackup = !a.ToggleShowBackup
			a.App.SetRoot(a.Root, true)
			a.App.SetFocus(a.VaultList)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			a.App.Stop()
			return nil
		}
		return event
	})

	return backupRoot, nil
}

func (a *App) BackupListView(backupDirEntries []os.DirEntry) *tview.Flex {
	backupList := tview.NewList()
	for _, v := range backupDirEntries {
		backupList.AddItem(v.Name(), "", 0, func() {})
	}
	backupList.SetBorder(true)
	backupList.SetTitle(" Backups ")
	backupList.SetBackgroundColor(tcell.ColorBlack)

	box := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(box, 0, 1, false).
		AddItem(backupList, 0, 1, true).
		AddItem(box, 0, 1, false)
}
