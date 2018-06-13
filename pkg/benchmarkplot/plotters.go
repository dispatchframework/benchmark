package benchmarkplot

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/onsi/ginkgo/types"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

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
