package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"regexp"
	"time"

	"github.com/dispatchframework/benchmark/pkg/reporter"
	. "github.com/logrusorgru/aurora"
)

var (
	testsRun  []string
	Functions []string
	// Flags to control how the tests are run
	samples  int
	testFunc string
	output   string
	apiIP    string
)

type Tester struct {
	name       string
	functions  []string
	apis       []string
	aggregator *reporter.BenchmarkRecorder
}

func init() {
	flag.StringVar(&output, "output",
		fmt.Sprintf("out-%v.csv", time.Now().Unix()),
		"What file to output the results to")
	flag.IntVar(&samples, "samples", 1, "Number of samples to be collected")
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
}

func (t *Tester) Cleanup() {
	fmt.Println("Cleaning up")
	fmt.Printf("Functions to be cleaned: %v\n", t.functions)
	for _, name := range t.functions {
		cmd := exec.Command("dispatch", "delete", "function", name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Unable to delete function %v. %v\n", name, err)
			log.Println("Can't delete function")
		}
	}
	fmt.Printf("Apis to be cleaned: %v\n", t.apis)
	for _, name := range t.apis {
		cmd := exec.Command("dispatch", "delete", "api", name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Unable to delete function %v. %v\n", name, err)
			log.Println("Can't delete function")
		}
	}

}

func callMethods(t *Tester, rx *regexp.Regexp) {
	for i := 0; i < reflect.ValueOf(t).NumMethod(); i++ {
		method := reflect.TypeOf(t).Method(i)
		name := method.Name
		if rx.MatchString(name) {
			fmt.Printf("\n\n[%v]\n\n", Green(name))
			reflect.ValueOf(t).MethodByName(name).Call(nil)
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	var testsMatcher string
	if len(args) > 0 {
		testsMatcher = args[0]
	} else {
		testsMatcher = "Test"
	}
	graphs := map[string]func(map[string][]float64, string){
		"Creation":  reporter.SeriesPlot,
		"Execution": reporter.SeriesPlot,
		"Scale":     reporter.SeriesPlot,
		"Api":       reporter.BarPlot,
	}
	rx := regexp.MustCompile(testsMatcher)
	testRecorder := reporter.NewReporter("TestTime", output)
	testRecorder.Graphs = graphs
	tests := &Tester{
		name:       "TimeTester",
		aggregator: testRecorder,
	}
	defer tests.Cleanup()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Signal Received: %v\n", sig)
			tests.Cleanup()
		}
	}()
	callMethods(tests, rx)
	fmt.Println(tests.aggregator.PrintResults())
}
