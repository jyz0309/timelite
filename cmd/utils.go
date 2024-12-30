package cmd

import (
	"github.com/rivo/tview"
)

func Error(err error) {
	page := getRootPage()
	page.RemovePage("error")
	modal := tview.NewModal().SetText(err.Error()).
		AddButtons([]string{"Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				page.HidePage("error")
			}
		})
	page.AddPage("error", modal, false, false)
	page.ShowPage("error")
}

func Success(msg string) {
	page := getRootPage()
	page.RemovePage(msg)

	modal := tview.NewModal().SetText(msg).
		AddButtons([]string{"Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				page.HidePage(msg)
			}
		})
	page.AddPage(msg, modal, false, false)
	page.ShowPage(msg)
}
