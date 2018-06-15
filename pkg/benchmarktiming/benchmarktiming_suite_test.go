package benchmarktiming_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/nickaashoek/benchmark/pkg/benchmarkreporter"
	. "github.com/nickaashoek/benchmark/pkg/benchmarktiming"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	wg         sync.WaitGroup
	workers    chan *Worker
	numWorkers int
	testFunc   string
	outputFile string
	shouldPlot bool
)

func init() {
	fmt.Println("Running Init")
	flag.StringVar(&outputFile, "OutputFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
	flag.BoolVar(&shouldPlot, "plot", false, "Should a plot be produced")
	fmt.Println("Ran Init")
}

func TestDispatch(t *testing.T) {
	RegisterFailHandler(Fail)
	reporter := NewDispatchReporter(outputFile, shouldPlot, nil)
	reporters := []Reporter{reporter}
	RunSpecsWithDefaultAndCustomReporters(t, "Dispatch Suite", reporters)
}

var _ = BeforeSuite(func() {
	numWorkers = 4
	workers = make(chan *Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		name := RandomName(10)
		newWorker := Worker{Me: i, Function: name}
		workers <- &newWorker
	}
	fmt.Println("Created workers")
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		wk := <-workers
		go func(location string, wk *Worker) {
			defer wg.Done()
			defer GinkgoRecover()
			createErr := wk.CreateFunction(location)
			Ω(createErr).Should(BeNil())
		}(testFunc, wk)
		workers <- wk
	}
	wg.Wait()
	fmt.Println("Done Setting up suite")
})

var _ = AfterSuite(func() {
	log.Printf("Running AfterSuite\n")
	for i := 0; i < numWorkers; i++ {
		wk := <-workers
		_ = wk.DeleteFunction()
	}
})

var _ = Describe("Measuring Function Creation Times", func() {

	var creationWorkers chan *Worker

	BeforeEach(func() {
		log.Println("Running BeforeEach")
		creationWorkers = make(chan *Worker, numWorkers)
		for i := 0; i < numWorkers; i++ {
			name := RandomName(10)
			newWorker := Worker{Me: i, Function: name}
			creationWorkers <- &newWorker
		}

	})

	AfterEach(func() {
		log.Println("Running AfterEach")
		DPrintf("Running AfterEach\n")
		for i := 0; i < numWorkers; i++ {
			wk := <-creationWorkers
			_ = wk.DeleteFunction()
		}
	})

	FMeasure("The time it takes to create a single function", func(b Benchmarker) {
		wk := <-creationWorkers
		createtime := b.Time("createtime", func() {
			createErr := wk.CreateFunction(testFunc)
			Ω(createErr).Should(BeNil())
		})
		creationWorkers <- wk
		DPrintf("Took %v seconds to create a single function\n", createtime)
	}, 1)

	Measure("Time it takes to create multiple functions in series", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < numWorkers; i++ {
				wk := <-creationWorkers
				createErr := wk.CreateFunction(testFunc)
				Ω(createErr).Should(BeNil())
				creationWorkers <- wk
			}
		})
		DPrintf("Run took %v second\n", runtime)
	}, 1)

	Measure("Time it takes to create functions in parallel", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				wk := <-creationWorkers
				go func(location string, wk *Worker) {
					defer wg.Done()
					defer GinkgoRecover()
					createErr := wk.CreateFunction(location)
					Ω(createErr).Should(BeNil())
				}(testFunc, wk)
				creationWorkers <- wk
			}
			wg.Wait()
			DPrintf("Finished\n")
		})
		DPrintf("Run took %v second\n", runtime)
	}, 1)
})

var _ = Describe("Measure execution time of functions", func() {

	Measure("The time it takes to run a mildly computationally intensive function", func(b Benchmarker) {
		wk := <-workers
		runTime := b.Time("runtime", func() {
			runErr := wk.ExecuteFunction(true)
			Ω(runErr).Should(BeNil())
		})
		workers <- wk
		DPrintf("Run took %v seconds\n", runTime)
	}, 1)

	Measure("Running the same function in series", func(b Benchmarker) {
		_ = b.Time("runtime", func() {
			for i := 0; i < numWorkers; i++ {
				wk := <-workers
				runErr := wk.ExecuteFunction(true)
				fmt.Println("Ran a single function")
				Ω(runErr).Should(BeNil())
				workers <- wk
			}
		})
	}, 1)

	Measure("Running the same function in parallel", func(b Benchmarker) {
		_ = b.Time("runtime", func() {
			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				wk := <-workers
				go func(wk *Worker) {
					defer wg.Done()
					defer GinkgoRecover()
					runErr := wk.ExecuteFunction(true)
					Ω(runErr).Should(BeNil())
				}(wk)
				workers <- wk
			}
			wg.Wait()
		})
	}, 1)

})
