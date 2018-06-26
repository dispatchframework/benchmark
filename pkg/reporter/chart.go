package reporter

import (
	"bufio"
	"os"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

func SimplePlot(Records map[string][]time.Duration, name string) {
	var series []chart.Series
	for field, durations := range Records {
		var x, y []float64
		for i, time := range durations {
			x = append(x, float64(i))
			y = append(y, time.Seconds())
		}
		series = append(series, chart.ContinuousSeries{
			Name:    field,
			XValues: x,
			YValues: y,
		})
	}
	PlotSeries(series, name)
}

func PlotSeries(series []chart.Series, name string) {
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
