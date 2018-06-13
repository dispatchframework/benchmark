package DispatchFunctionScaling_test

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/nickaashoek/benchmark/pkg/DispatchFunctionTiming"
	. "github.com/nickaashoek/benchmark/pkg/dispatch-reporter"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	wg          sync.WaitGroup
	testFunc    string
	outputFile  string
	runFunction *exec.Cmd
)

func MkPlot(measurements map[string]types.SpecMeasurement) {
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

func parallelExec(runs int, command *exec.Cmd) {
	for i := 0; i < runs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			toRun := *command
			err := toRun.Run()
			if err != nil {
				log.Fatalf("Function failed to executed: %v\n", err)
			}
		}()
	}
	wg.Wait()
}

func init() {
	flag.StringVar(&outputFile, "outFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
}

func TestDispatchFunctionScaling(t *testing.T) {
	RegisterFailHandler(Fail)
	// outputFile = fmt.Sprintf("./output-%v.csv", time.Now().Unix())
	reporter := NewDispatchReporter(outputFile, true, MkPlot)
	reporters := []Reporter{reporter}
	RunSpecsWithDefaultAndCustomReporters(t, "Dispatch Suite", reporters)
}

var _ = BeforeSuite(func() {
	testFunc = "../../resources/functions/test.py"
	runFunction = exec.Command("dispatch", "exec", "scaling-test", "--wait")
	createWorker := Worker{Me: 0, Function: "scaling-test"}
	err := createWorker.CreateFunction(testFunc)
	if err != nil {
		log.Fatal("Failed to create the function")
	}
})

var _ = AfterSuite(func() {
	teardown := exec.Command("dispatch", "delete", "function", "scaling-test")
	err := teardown.Run()
	if err != nil {
		log.Fatal("Unable to teardown the function")
	}
})

var _ = Describe("Testing a simple functions run at different scales", func() {
	var (
		command exec.Cmd
	)
	BeforeEach(func() {
		command = *runFunction
	})
	MeasureRuntime := func(execs int) {
		Measure(fmt.Sprintf("%v runs in parallel", execs), func(b Benchmarker) {
			_ = b.Time(fmt.Sprintf("%v", execs), func() {
				parallelExec(execs, runFunction)
			})
		}, 5)
	}
	Context("Running the function in parallel sets", func() {
		for i := 1; i < 2048; i *= 2 {
			MeasureRuntime(i)
		}
	})
})
