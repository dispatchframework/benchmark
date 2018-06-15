package benchmarkplot

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/onsi/ginkgo/types"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func ApiThroughputPlot(measurements map[string]types.SpecMeasurement) {
	p, err := plot.New()
	if err != nil {
		log.Fatal("Unable to create plot")
	}
	p.Title.Text = "Results of throughput experiment"
	p.X.Label.Text = "Number of queriers"
	p.Y.Label.Text = "Number of successful API Queries in 1 Second"

	var pts plotter.XYs
	for field, measurement := range measurements {
		if strings.Contains(field, "queriers") {
			queriers, _ := strconv.ParseFloat(measurement.Name, 64)
			pts = append(pts, struct {
				X float64
				Y float64
			}{queriers, measurement.Average})
		}
	}
	sort.Slice(pts[:], func(i, j int) bool {
		return pts[i].X < pts[j].X
	})
	err = plotutil.AddLinePoints(p, "Scalability", pts)
	if err != nil {
		log.Fatal("Failed to add points")
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func ScalePlot(measurements map[string]types.SpecMeasurement) {
	p, err := plot.New()
	if err != nil {
		log.Fatal("Unable to create plot")
	}
	p.Title.Text = "Results of function scaling measurements"
	p.X.Label.Text = "Number of functions run in parallel"
	p.Y.Label.Text = "Time (s)"

	var pts plotter.XYs
	for _, measurement := range measurements {
		execs, _ := strconv.ParseFloat(measurement.Name, 64)
		fmt.Printf("Execs: %v, %s\n", execs, measurement.Name)
		pts = append(pts, struct {
			X float64
			Y float64
		}{execs, measurement.Average})
	}
	sort.Slice(pts[:], func(i, j int) bool {
		return pts[i].X < pts[j].X
	})
	err = plotutil.AddLinePoints(p, "Scalability", pts)
	if err != nil {
		log.Fatal("Failed to add points")
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}
