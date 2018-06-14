package benchmarkscaling_test

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"testing"
	"time"

	. "github.com/nickaashoek/benchmark/pkg/benchmarkplot"
	. "github.com/nickaashoek/benchmark/pkg/benchmarkreporter"
	. "github.com/nickaashoek/benchmark/pkg/benchmarktiming"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	wg          sync.WaitGroup
	testFunc    string
	outputFile  string
	runFunction *exec.Cmd
)

func init() {
	flag.StringVar(&outputFile, "outFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
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

func TestDispatchFunctionScaling(t *testing.T) {
	RegisterFailHandler(Fail)
	// outputFile = fmt.Sprintf("./output-%v.csv", time.Now().Unix())
	reporter := NewDispatchReporter(outputFile, true, ScalePlot)
	reporters := []Reporter{reporter}
	RunSpecsWithDefaultAndCustomReporters(t, "Dispatch Suite", reporters)
}

var _ = BeforeSuite(func() {
	testFunc = "/Users/nkaashoek/go/src/github.com/nickaashoek/benchmark/resources/functions/test.py"
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
		for i := 1; i < 4; i *= 2 {
			MeasureRuntime(i)
		}
	})
})
