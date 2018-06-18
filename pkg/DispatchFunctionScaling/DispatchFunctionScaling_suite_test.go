package DispatchFunctionScaling_test

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"testing"
	"time"

	. "github.com/dispatchframework/benchmark/pkg/DispatchFunctionTiming"
	. "github.com/dispatchframework/benchmark/pkg/dispatch-reporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	wg          sync.WaitGroup
	testFunc    string
	outputFile  string
	runFunction *exec.Cmd
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

func init() {
	flag.StringVar(&outputFile, "outFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
}

func TestDispatchFunctionScaling(t *testing.T) {
	RegisterFailHandler(Fail)
	// outputFile = fmt.Sprintf("./output-%v.csv", time.Now().Unix())
	reporter := NewDispatchReporter(outputFile)
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
	var command exec.Cmd
	BeforeEach(func() {
		command = *runFunction
	})
	Measure("Running the function once as a baseline", func(b Benchmarker) {
		runtime := b.Time("Single run of function", func() {
			err := command.Run()
			if err != nil {
				log.Fatalf("%s\n", err)
			}
			// Ω(err).Should(BeNil())
		})
		Ω(runtime.Seconds()).Should(BeNumerically("<", 3), "Function is simple, shouldn't take more than a few seconds")
	}, 1)

	Measure("Running the function four times in parallel", func(b Benchmarker) {
		_ = b.Time("Parallel runs", func() {
			parallelExec(4, runFunction)
		})
	}, 1)

	Measure("Running the function 16 times in parallel", func(b Benchmarker) {
		_ = b.Time("Parallel runs", func() {
			parallelExec(16, runFunction)
		})
	}, 1)

	Measure("Running the function 64 times in parallel", func(b Benchmarker) {
		_ = b.Time("Parallel runs", func() {
			parallelExec(64, runFunction)
		})
	}, 1)

	Measure("Running the function 128 times in parallel", func(b Benchmarker) {
		_ = b.Time("Parallel runs", func() {
			parallelExec(128, runFunction)
		})
	}, 1)

	Measure("Running the function 512 times in parallel", func(b Benchmarker) {
		_ = b.Time("Parallel runs", func() {
			parallelExec(512, runFunction)
		})
	}, 1)

})
