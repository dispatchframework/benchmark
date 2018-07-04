package reporter

import (
	"bufio"
	"fmt"
	"os"

	chart "github.com/wcharczuk/go-chart"
)

// SeriesPlot formats the records correctly for go-chart series plot to handle
func SeriesPlot(Records map[string][]float64, name string) {
	var series []chart.Series
	for field, samples := range Records {
		var x, y []float64
		for i, record := range samples {
			x = append(x, float64(i))
			y = append(y, record)
		}
		series = append(series, chart.ContinuousSeries{
			Name:    field,
			XValues: x,
			YValues: y,
		})
	}
	name = fmt.Sprintf("%v-chart.png", name)
	plotSeries(series, name)
}

func plotSeries(series []chart.Series, name string) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Series: series,
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	var output *os.File
	defer output.Close()
	if _, err := os.Stat(name); err == nil {
		os.Remove(name)
	}
	output, _ = os.Create(name)
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	_ = graph.Render(chart.PNG, writer)
}

// BarPlot formats the records correctly for go-chart bar plot to handle
func BarPlot(Records map[string][]float64, name string) {
	var bars []chart.Value
	for field, samples := range Records {
		value, _ := GetStats(samples)
		bars = append(bars, chart.Value{
			Value: value,
			Label: field,
		})
	}
	name = fmt.Sprintf("%v-chart.png", name)
	plotBar(bars, name)
}

func plotBar(bars []chart.Value, name string) {
	graph := chart.BarChart{
		Title:      name,
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 30,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: bars,
	}
	var output *os.File
	defer output.Close()
	if _, err := os.Stat(name); err == nil {
		os.Remove(name)
	}
	output, _ = os.Create(name)
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	_ = graph.Render(chart.PNG, writer)

}
