package cmd

import (
	"timelite/conf"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

var panelDetailPage *PanelDetailPage

func getPanelDetailPage(conf *conf.Panel, dashboard *conf.Dashboard) *PanelDetailPage {
	if panelDetailPage == nil {
		form := tview.NewForm().
			AddInputField("Title", "", 100, nil, nil).
			AddInputField("Query", "", 100, nil, nil).
			AddDropDown("Panel Type", []string{"plot", "gauge"}, 0, nil).
			AddButton("Update", func() {
				panelDetailPage.Update()
			}).
			AddButton("Delete", func() {
				panelDetailPage.Delete()
			}).
			AddButton("Back", func() {
				getRootPage().SwitchToPage("dashboards")
			}).
			AddButton("Refresh", func() {
				panelDetailPage.Refresh()
			})
		form.SetBorder(true)

		panelView := tview.NewFlex().SetDirection(tview.FlexRow)
		panelView.SetBorder(true)

		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.AddItem(panelView, 0, 3, false)
		flex.AddItem(form, 0, 1, true)

		panelDetailPage = &PanelDetailPage{
			panelView: panelView,
			panelForm: form,
			flex:      flex,
		}
	}
	panelDetailPage.panel = conf
	panelDetailPage.dashboard = dashboard
	return panelDetailPage
}

type PanelDetailPage struct {
	panel     *conf.Panel
	dashboard *conf.Dashboard

	panelView *tview.Flex
	panelForm *tview.Form

	flex *tview.Flex
}

func (page *PanelDetailPage) Visualize() tview.Primitive {
	page.panelForm.GetFormItemByLabel("Title").(*tview.InputField).SetText(page.panel.Title)
	page.panelForm.GetFormItemByLabel("Query").(*tview.InputField).SetText(page.panel.Query)

	if page.panelView.GetItemCount() > 0 {
		// remove old panel
		page.panelView.RemoveItem(page.panelView.GetItem(0))
	}

	data, ts, err := query(page.panel.Query)
	if err != nil {
		logrus.Errorf("failed to query promql, err[%s]", err.Error())
		dialog := tvxwidgets.NewMessageDialog(tvxwidgets.ErrorDailog)
		dialog.SetMessage(err.Error())
		page.panelView.AddItem(dialog, 0, 3, false)
	} else {
		page.panelView.AddItem(createPanel(page.panel, data, ts), 0, 3, false)
	}
	return page.flex
}

func (page *PanelDetailPage) Refresh() {
	data, ts, err := query(page.panel.Query)
	if err != nil {
		logrus.Errorf("failed to query promql, err[%s]", err.Error())
		Error(err)
	} else {
		page.panelView.RemoveItem(page.panelView.GetItem(0))
		page.panelView.SetTitle(page.panel.Title)
		page.panelView.AddItem(createPanel(page.panel, data, ts), 0, 3, false)
	}
}

func (page *PanelDetailPage) Delete() {
	if err := page.dashboard.DeletePanel(page.panel); err != nil {
		logrus.Errorf("failed to delete panel, err[%s]", err.Error())
		Error(err)
	} else {
		Success("Delete panel successfully!")
		// refresh dashboard page after delete
		getDashboardPage().RefreshDashboard()

		getRootPage().SwitchToPage("dashboards")
	}
}

func (page *PanelDetailPage) Update() {
	_, panelType := page.panelForm.GetFormItemByLabel("Panel Type").(*tview.DropDown).GetCurrentOption()

	newPanel := &conf.Panel{
		Title:     page.panelForm.GetFormItemByLabel("Title").(*tview.InputField).GetText(),
		Query:     page.panelForm.GetFormItemByLabel("Query").(*tview.InputField).GetText(),
		PanelType: panelType,
	}
	if err := page.dashboard.UpdatePanel(page.panel, newPanel); err != nil {
		logrus.Errorf("failed to update panel, err[%s]", err.Error())
		Error(err)
	} else {
		Success("Update panel successfully!")
	}
	page.panel = newPanel

	// refresh dashboard page after update
	page.Refresh()
	getDashboardPage().RefreshDashboard()
}
