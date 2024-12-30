package cmd

import (
	"context"
	"time"
	"timelite/conf"
	query_engine "timelite/query"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

func QueryPage() tview.Primitive {
	page := getRootPage()
	form := tview.NewForm().AddTextArea("PromQL", "Enter promql here...", 100, 5, 0, nil).
		AddDropDown("Panel", []string{"plot", "gauge", "bar"}, 0, nil)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(form, 0, 1, false)
	flex.AddItem(nil, 0, 3, false)

	form.AddButton("Query", func() {
		logrus.Info("press the query button")
		textArea := form.GetFormItemByLabel("PromQL").(*tview.TextArea)
		promql := textArea.GetText()
		_, panelType := form.GetFormItemByLabel("Panel").(*tview.DropDown).GetCurrentOption()
		data, ts, err := query(promql)
		if err != nil {
			logrus.Errorf("failed to query promql, err[%s]", err.Error())
			Error(err)
			return
		}
		panel := createPanel(&conf.Panel{
			Title:     "Query",
			PanelType: panelType,
			Query:     promql,
		}, data, ts)
		panel.SetBorder(true)
		panel.SetTitle("Query")

		flex.RemoveItem(flex.GetItem(1))
		flex.AddItem(panel, 0, 3, false)
	}).
		AddButton("Save", func() {
			textArea := form.GetFormItem(0).(*tview.TextArea)
			promql := textArea.GetText()
			_, panelType := form.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
			if !page.HasPage("save_modal") {
				page.AddPage("save_modal", saveModal(promql, panelType), true, false)
			}
			page.ShowPage("save_modal")
		}).
		AddButton("Back", func() {
			page.SwitchToPage("main")
		})

	form.SetBorder(true).
		SetTitle("Query").
		SetTitleAlign(tview.AlignLeft)

	flex.SetTitle("Query").
		SetTitleAlign(tview.AlignCenter)
	flex.SetFullScreen(true)

	return flex
}

func query(qs string) ([]*query_engine.Series, []int64, error) {
	now := time.Now()
	querier := query_engine.GetGlobalQuerier(conf.DefaultConfig.StoragePath)
	return querier.NewRangeQuery(
		context.Background(),
		qs,
		now.Truncate(15*time.Minute),
		now,
		15*time.Second)
}

func saveModal(promql string, panelType string) tview.Primitive {
	dashboards := GetDashboardsTree()

	window := tview.NewFlex().SetDirection(tview.FlexRow)

	form := tview.NewForm().
		AddInputField("title", "", 0, nil, nil)

	form.AddButton("Save", func() {
		if ref := dashboards.GetCurrentNode().GetReference(); ref != nil {
			title := form.GetFormItemByLabel("title").(*tview.InputField)
			dashboard := ref.(*conf.Dashboard)
			err := save(title.GetText(), panelType, promql, dashboard)
			if err != nil {
				Error(err)
			} else {
				Success("Save panel to dashboard successfully!")
			}
		}
	}).
		AddButton("Back", func() {
			page := getRootPage()
			page.RemovePage("save_modal")
		})

	window.AddItem(dashboards, 0, 3, false)
	window.AddItem(form, 5, 0, false).
		SetTitle("Save which dashboard?").
		SetTitleAlign(tview.AlignLeft)
	window.SetBorder(true)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 2, false).
			AddItem(window, 0, 1, true).
			AddItem(nil, 0, 2, false), 0, 1, true).
		AddItem(nil, 0, 1, false)
	return modal
}

func save(title, panelType, query string, dashboard *conf.Dashboard) error {
	err := dashboard.AddPanel(&conf.Panel{
		Title:     title,
		PanelType: panelType,
		Query:     query,
	})
	if err != nil {
		return err
	}
	err = dashboard.Save()
	if err != nil {
		return err
	}
	return nil
}
