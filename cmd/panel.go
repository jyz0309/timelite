package cmd

import (
	"fmt"
	"time"
	"timelite/conf"
	query_engine "timelite/query"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

var colors = []tcell.Color{
	tcell.ColorSteelBlue,
	tcell.ColorGreen,
	tcell.ColorRed,
	tcell.ColorYellow,
	tcell.ColorBlue,
	tcell.ColorPurple,
	tcell.ColorOrange,
	tcell.ColorPink,
	tcell.ColorBrown,
	tcell.ColorGray,
	tcell.ColorBlack,
	tcell.ColorWhite,
	tcell.ColorAqua,
	tcell.ColorFuchsia,
	tcell.ColorTeal,
	tcell.ColorOlive,
	tcell.ColorMaroon,
	tcell.ColorNavy,
	tcell.ColorLime,
	tcell.ColorSilver,
	tcell.ColorBeige,
	tcell.ColorBisque,
	tcell.ColorBlanchedAlmond,
	tcell.ColorBlueViolet,
	tcell.ColorBrown,
	tcell.ColorCadetBlue,
	tcell.ColorChartreuse,
	tcell.ColorChocolate,
	tcell.ColorCoral,
	tcell.ColorCornflowerBlue,
	tcell.ColorCornsilk,
	tcell.ColorCrimson,
	tcell.ColorDarkBlue,
	tcell.ColorDarkCyan,
	tcell.ColorDarkGoldenrod,
	tcell.ColorDarkGray,
	tcell.ColorDarkGreen,
	tcell.ColorDarkKhaki,
	tcell.ColorDarkMagenta,
	tcell.ColorDarkOliveGreen,
	tcell.ColorDarkOrange,
	tcell.ColorDarkOrchid,
	tcell.ColorDarkRed,
	tcell.ColorDarkSalmon,
	tcell.ColorDarkSeaGreen,
	tcell.ColorDarkSlateBlue,
	tcell.ColorDarkSlateGray,
	tcell.ColorDarkTurquoise,
	tcell.ColorDarkViolet,
	tcell.ColorDeepPink,
}

func newPlot(title string, series []*query_engine.Series, ts []int64) *tvxwidgets.Plot {
	data := make([][]float64, len(series))
	for i, serie := range series {
		data[i] = serie.Floats
	}
	plot := tvxwidgets.NewPlot()
	plot.SetTitle(title)
	plot.SetData(data)
	plot.SetMarker(tvxwidgets.PlotMarkerBraille)
	plot.SetXAxisLabelFunc(func(i int) string {
		return time.UnixMilli(ts[i]).Format("15:04:05")
	})

	plot.SetLineColor(colors)
	plot.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDoubleClick {
			if plot.HasFocus() {
				logrus.Infof("press the plot[%s]", plot.GetTitle())
			}
		}
		return action, event
	})
	return plot
}

func newGauge(series []*query_engine.Series) *tview.Flex {
	gaugeFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, v := range series {
		if len(v.Floats) > 0 {
			gauge := tvxwidgets.NewUtilModeGauge()
			gauge.SetValue(v.Floats[len(v.Floats)-1])
			gauge.SetLabel(fmt.Sprintf("%s", v.Metric))
			gaugeFlex.AddItem(gauge, 1, 0, false)
		}
	}
	return gaugeFlex
}

// TODO: Now the bar chart is not working well, it only allows int value,
// but the data is always float64, so we need to convert it to int first.
// Will fix it later by modifying the tvxwidgets.
func newBarChart(title string, series []*query_engine.Series) *tvxwidgets.BarChart {
	barChart := tvxwidgets.NewBarChart()
	barChart.SetTitle(title)
	for i, v := range series {
		barChart.AddBar(fmt.Sprintf("%s", v.Metric), int(v.Floats[len(v.Floats)-1]), colors[i%len(colors)])
	}
	return barChart
}

func createPanel(p *conf.Panel, series []*query_engine.Series, ts []int64) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	var panel tview.Primitive
	switch p.PanelType {
	case "plot":
		plot := newPlot(p.Title, series, ts)
		panel = plot
	case "gauge":
		gauge := newGauge(series)
		panel = gauge
	case "bar":
		barChart := newBarChart(p.Title, series)
		panel = barChart
	}
	flex.AddItem(panel, 0, 1, true)
	return flex
}
