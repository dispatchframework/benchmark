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
	. "github.com/logrusorgru/aurora"
)

var (
	testsRun  []string
	Functions []string
	// Flags to control how the tests are run
	samples    int
	testFunc   string
	shouldPlot bool
	output     string
	apiIP      string
)

type Tester interface {
	Cleanup()
}

type TimeTests struct {
	name       string
	functions  []string
	apis       []string
	aggregator *reporter.TimeRecord
}

type ScaleTests struct {
	name       string
	functions  []string
	aggregator *reporter.TimeRecord
}

func init() {
	flag.StringVar(&output, "outFile",
		fmt.Sprintf("out-%v.csv", time.Now().Unix()),
		"What file to output the results to")
	flag.IntVar(&samples, "samples", 1, "Number of samples to be collected")
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
	flag.BoolVar(&shouldPlot, "plot", true, "Should a plot be produced")
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
	fmt.Printf("Apis to be cleaned: %v\n", t.apis)
	for _, name := range t.apis {
		cmd := exec.Command("dispatch", "delete", "api", name)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Unable to delete function %v. %v\n", name, err)
			log.Println("Can't delete function")
		}
	}

}

func (t *ScaleTests) Cleanup() {
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

func callMethods(t Tester, rx *regexp.Regexp) {
	for i := 0; i < reflect.ValueOf(t).NumMethod(); i++ {
		method := reflect.TypeOf(t).Method(i)
		name := method.Name
		if rx.MatchString(name) {
			fmt.Printf("\n\n[%v]\n\n", Green(name))
			reflect.ValueOf(t).MethodByName(name).Call(nil)
		}
	}
	t.Cleanup()
}

func main() {
	flag.Parse()
	args := flag.Args()
	var testsMatcher string
	if len(args) > 0 {
		testsMatcher = args[0]
	} else {
		testsMatcher = ".+"
	}
	rx := regexp.MustCompile(testsMatcher)
	timeRecorder := reporter.NewReporter("TestTime", output, shouldPlot, reporter.SimplePlot)
	scaleRecorder := reporter.NewReporter("TestScale", output, shouldPlot, reporter.SimplePlot)
	timer := &TimeTests{
		name:       "TimeTester",
		aggregator: timeRecorder,
	}
	scales := &ScaleTests{
		name:       "ScaleTester",
		aggregator: scaleRecorder,
	}
	callMethods(timer, rx)
	callMethods(scales, rx)
	fmt.Println(timer.aggregator.PrintResults())
	fmt.Println(scales.aggregator.PrintResults())
}
