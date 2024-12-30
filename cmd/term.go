package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"timelite/conf"
	query_engine "timelite/query"
)

var rootPage *tview.Pages

func getRootPage() *tview.Pages {
	if rootPage == nil {
		rootPage = tview.NewPages()
	}
	return rootPage
}

func Init() {
	app := tview.NewApplication().EnableMouse(true)
	rootPage := getRootPage()

	rootPage.AddPage("query", QueryPage(), true, false)
	rootPage.AddPage("dashboards", getDashboardPage().Visualize(), true, false)
	list := tview.NewList().
		AddItem("Query", "", 'a', func() {
			rootPage.SwitchToPage("query")
		}).
		AddItem("Dashboards", "", 'b', func() {
			rootPage.SwitchToPage("dashboards")
		}).
		AddItem("Quit", "", 'q', func() {
			query_engine.GetGlobalQuerier(conf.DefaultConfig.StoragePath).Close()
			app.Stop()
		})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			query_engine.GetGlobalQuerier(conf.DefaultConfig.StoragePath).Close()
			app.Stop()
		}
		return event
	})
	list.ShowSecondaryText(false)
	rootPage.AddPage("main", list, true, true)
	if err := app.SetRoot(rootPage, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

}
