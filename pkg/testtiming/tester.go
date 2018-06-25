package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
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
	name      string
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

func (t *TimeTests) Cleanup() {
	fmt.Println("Cleaning up")
	fmt.Printf("Functions to be cleaned: %v\n", t.functions)
	for _, name := range t.functions {
		cmd := exec.Command("dispatch", "delete", "function", name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Unable to delete function %v. %v\n", name, err)
			log.Println("Can't delete function")
		}
	}
}

func main() {
	flag.Parse()
	testsMatcher := os.Args[1]
	rx := regexp.MustCompile(testsMatcher)
	aggregator = reporter.NewReporter("test runner", output, shouldPlot)
	tests := &TimeTests{
		name: "Testing",
	}
	for i := 0; i < reflect.ValueOf(tests).NumMethod(); i++ {
		method := reflect.TypeOf(tests).Method(i)
		name := method.Name
		if rx.MatchString(name) {
			reflect.ValueOf(tests).MethodByName(name).Call(nil)
		}
	}
	tests.Cleanup()
	// tests.TestFuncMake()
	// tests.TestFuncMakeSerial()
	fmt.Println(aggregator.PrintResults())
}
