package benchmarkscaling_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	. "github.com/dispatchframework/benchmark/pkg/benchmarkplot"
	. "github.com/dispatchframework/benchmark/pkg/benchmarkreporter"
	. "github.com/dispatchframework/benchmark/pkg/benchmarktiming"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	wg          sync.WaitGroup
	testFunc    string
	outputFile  string
	runFunction *exec.Cmd
	shouldPlot  bool
)

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

func TestDispatchFunctionScaling(t *testing.T) {
	RegisterFailHandler(Fail)
	reporter := NewDispatchReporter(outputFile, shouldPlot, ScalePlot)
	reporters := []Reporter{reporter}
	RunSpecsWithDefaultAndCustomReporters(t, "Dispatch Suite", reporters)
}

var _ = BeforeSuite(func() {
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

var _ = Describe("", func() {

	var (
		command exec.Cmd
		samples int
	)
	fmt.Printf("Samples: %v\n", samples)
	flag.IntVar(&samples, "samples", 1, "Number of samples to be collected")
	flag.StringVar(&outputFile, "outFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
	fmt.Printf("OUTPUT FILE: %v\n", outputFile)
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
	flag.BoolVar(&shouldPlot, "plot", false, "Should a plot be produced")
	flag.Parse()
	Context("Testing a simple functions run at different scales", func() {
		BeforeEach(func() {
			command = *runFunction
		})
		Context("Running the function in parallel sets", func() {
			for i := 1; i < 2; i *= 2 {
				Measure(fmt.Sprintf("%v runs in parallel", i), func(b Benchmarker) {
					_ = b.Time(fmt.Sprintf("%v", i), func() {
						parallelExec(i, runFunction)
					})
				}, samples)
			}
		})
	})
})
