package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/dispatchframework/benchmark/pkg/reporter"
)

var (
	aggregator *reporter.TimeRecord
	testsRun   []string
	Functions  []string
	// Flags to control how the tests are run
	samples    int
	testFunc   string
	shouldPlot bool
	output     string
)

type TimeTests struct {
	functions []string
	apis      []string
}

func init() {
	flag.StringVar(&output, "outFile",
		fmt.Sprintf("out-%v.csv", time.Now().Unix()),
		"What file to output the results to")
	flag.IntVar(&samples, "samples", 1, "Number of samples to be collected")
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
	flag.BoolVar(&shouldPlot, "plot", false, "Should a plot be produced")
}

func (t TimeTests) Cleanup() {
	fmt.Println("Cleaning up")
	fmt.Printf("Functions: %v\n", Functions)
	for _, name := range Functions {
		cmd := exec.Command("dispatch", "delete", "function", name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Unable to delete function %v. %v\n", name, err)
			log.Println("Can't delete function")
		}
	}
}

func main() {
	flag.Parse()
	aggregator = reporter.NewReporter("test runner", output, shouldPlot)
	var tests TimeTests
	defer tests.Cleanup()
	tests.TestFuncMake()
	fmt.Println(aggregator.PrintResults())
}
