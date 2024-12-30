package cmd

import (
	"fmt"
	"timelite/conf"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type DashboardPage struct {
	visualizeFlex *tview.Flex
	listFlex      *tview.Flex
	panelFlex     *tview.Flex

	currentDashboard *conf.Dashboard
}

var dashboardPage *DashboardPage

func getDashboardPage() *DashboardPage {
	if dashboardPage != nil {
		return dashboardPage
	}
	listFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	listFlex.SetBorder(true)

	panelFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	panelFlex.SetBorder(true)

	dashboardPage = &DashboardPage{
		visualizeFlex: tview.NewFlex().SetDirection(tview.FlexColumn),
		listFlex:      listFlex,
		panelFlex:     panelFlex,
	}
	return dashboardPage
}

func (page *DashboardPage) Visualize() tview.Primitive {
	page.listFlex = page.ListDashboards()
	page.visualizeFlex.AddItem(page.listFlex, 0, 1, false)
	page.visualizeFlex.AddItem(page.panelFlex, 0, 4, false)
	page.visualizeFlex.SetFullScreen(true)
	return page.visualizeFlex
}

func (page *DashboardPage) RefreshDashboard() {
	logrus.Infof("Refresh dashboard[%s]", page.currentDashboard.Name)
	page.panelFlex = page.panelFlex.Clear()
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	panels := page.currentDashboard.GetAllPanels()
	count := 0
	for _, p := range panels {
		panel := page.createPanel(p)
		// 2 panels per row
		if count%2 == 0 {
			leftFlex.AddItem(panel, 0, 1, false)
		} else {
			rightFlex.AddItem(panel, 0, 1, false)
		}
		count++
	}

	// Fill leftFlex
	if count := leftFlex.GetItemCount(); count < 3 {
		for i := 0; i < 3-count; i++ {
			leftFlex.AddItem(nil, 0, 1, false)
		}
	}

	// Fill rightFlex
	if count := rightFlex.GetItemCount(); count < 3 {
		for i := 0; i < 3-count; i++ {
			rightFlex.AddItem(nil, 0, 1, false)
		}
	}

	page.panelFlex.
		AddItem(leftFlex, 0, 1, false).
		AddItem(rightFlex, 0, 1, false)
}

func (page *DashboardPage) ListDashboards() *tview.Flex {
	// if listFlex is empty, add button
	if page.listFlex.GetItemCount() == 0 {
		tree := GetDashboardsTree().SetSelectedFunc(func(node *tview.TreeNode) {
			if node.GetReference() == nil {
				return
			}
			dashboard := node.GetReference().(*conf.Dashboard)
			page.currentDashboard = dashboard
			page.RefreshDashboard()
		})

		rootPage := getRootPage()
		form := tview.NewForm().
			AddButton("Create", func() {
				if !rootPage.HasPage("create_dashboard_modal") {
					rootPage.AddPage("create_dashboard_modal", page.createDashboard(), true, false)
				}
				rootPage.ShowPage("create_dashboard_modal")
			}).
			AddButton("Delete", func() {
				if !rootPage.HasPage("delete_dashboard_modal") {
					rootPage.AddPage("delete_dashboard_modal", page.deleteDashboard(tree.GetCurrentNode().GetText()), true, false)
				}
				rootPage.ShowPage("delete_dashboard_modal")
			}).
			AddButton("Back", func() {
				rootPage.SwitchToPage("main")
			})
		form.SetButtonBackgroundColor(tcell.ColorBlack)

		page.listFlex.AddItem(form, 3, 0, false)
		page.listFlex.AddItem(tree, 0, 8, false)
		page.listFlex.SetTitle("Dashboards")
	}
	return page.listFlex
}

func (page *DashboardPage) deleteDashboard(name string) tview.Primitive {
	if name == "" || name == "." {
		return nil
	}
	text := fmt.Sprintf("[yellow]Are you sure to delete dashboard[red][%s][yellow]?", name)
	rootPage := getRootPage()

	form := tview.NewForm().
		AddTextView("", text, len(text), 3, true, false).
		AddButton("Delete", func() {
			if err := conf.DeleteDashboard(name); err != nil {
				Error(err)
			} else {
				tree := page.listFlex.GetItem(1).(*tview.TreeView)
				tree.GetRoot().RemoveChild(tree.GetCurrentNode())
				if page.currentDashboard.Name == name {
					page.currentDashboard = nil
				}
				Success("Delete dashboard successfully!")
				rootPage.HidePage("delete_dashboard_modal")
			}
		}).
		AddButton("Cancel", func() {
			rootPage := getRootPage()
			rootPage.RemovePage("delete_dashboard_modal")
		})
	form.SetBorder(true)
	form.SetTitle("Delete Dashboard")
	form.SetTitleAlign(tview.AlignLeft)
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 2, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 2, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	return flex
}

func (page *DashboardPage) createDashboard() tview.Primitive {
	form := tview.NewForm()
	rootPage := getRootPage()
	form.AddInputField("title", "", 20, nil, nil).
		AddButton("Create", func() {
			title := form.GetFormItemByLabel("title").(*tview.InputField)
			dashboard, err := conf.NewDashboard(title.GetText())
			if err != nil {
				Error(err)
			} else {
				tree := page.listFlex.GetItem(1).(*tview.TreeView)
				tree.GetRoot().AddChild(tview.NewTreeNode(dashboard.Name).
					SetReference(dashboard).
					SetSelectable(true).
					SetColor(tcell.ColorGreen))
				Success("Create dashboard successfully!")
				rootPage.HidePage("create_dashboard_modal")
			}
		}).
		AddButton("Cancel", func() {
			rootPage.HidePage("create_dashboard_modal")
		})
	form.SetBorder(true)
	form.SetTitle("Create Dashboard")
	form.SetTitleAlign(tview.AlignLeft)
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 2, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 2, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	return flex
}

func (page *DashboardPage) createPanel(p *conf.Panel) tview.Primitive {
	data, ts, err := query(p.Query)
	if err != nil {
		logrus.Errorf("failed to query promql, err[%s]", err.Error())
		dialog := tvxwidgets.NewMessageDialog(tvxwidgets.ErrorDailog)
		dialog.SetMessage(err.Error())
		return dialog
	}
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	rootPage := getRootPage()
	form := tview.NewForm().
		AddButton("Detail", func() {
			if !rootPage.HasPage("panel_detail") {
				rootPage.AddPage("panel_detail", getPanelDetailPage(p, page.currentDashboard).Visualize(), true, false)
			}
			// refresh panel detail page
			getPanelDetailPage(p, page.currentDashboard).Visualize()
			rootPage.SwitchToPage("panel_detail")
		}).
		AddButton("Delete", func() {
			page.currentDashboard.DeletePanel(p)
			page.RefreshDashboard()
		}).
		AddButton("Refresh", func() {
			data, ts, err := query(p.Query)
			if err != nil {
				logrus.Errorf("failed to query promql, err[%s]", err.Error())
				dialog := tvxwidgets.NewMessageDialog(tvxwidgets.ErrorDailog)
				dialog.SetMessage(err.Error())
				flex.Clear()
				flex.AddItem(dialog, 0, 1, true)
			} else {
				flex.RemoveItem(flex.GetItem(1))
				flex.AddItem(createPanel(p, data, ts), 0, 6, true)
			}
		}).SetButtonsAlign(tview.AlignRight)
	flex.AddItem(form, 0, 1, true).
		AddItem(createPanel(p, data, ts), 0, 6, true)
	flex.SetBorder(true)
	flex.SetTitle(p.Title)
	return flex
}

func GetDashboardsTree() *tview.TreeView {
	dashboards := conf.ListDashboards()

	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	for _, dashboard := range dashboards {
		node := tview.NewTreeNode(dashboard.Name).
			SetReference(dashboard).
			SetColor(tcell.ColorGreen).
			SetSelectable(true)
		root.AddChild(node)
	}
	return tree
}
