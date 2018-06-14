package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]
	config := ReadJson(args[0])
	if config.Timing.ToRun {
		fmt.Println("Run timing test")
		runTiming(config.Timing)
	}
	if config.Scaling.ToRun {
		fmt.Println("Run scaling test")
		runScaling(config.Scaling)
	}

}

func runScaling(config TestConfig) {
	fmt.Printf("Running Scalability Tests: %v\n", config)
	var testExec *exec.Cmd
	test := fmt.Sprintf("%v/benchmarkscale.test", config.Location)
	shouldPlot := fmt.Sprintf("-plot=%v", config.Plot)
	if config.Output == "" {
		testExec = exec.Command(test, shouldPlot)
	} else {
		output := fmt.Sprintf("-outFile=%v/%v", os.Getenv("PWD"), config.Output)
		testExec = exec.Command(test, output, shouldPlot)
	}
	output, err := testExec.Output()

	// Exit status 197 is a special error status used by ginkgo to reflect programatic focus,
	// we don't want to report a test as failed in this case
	if err != nil && err.Error() != "exit status 197" {
		fmt.Printf("Error: %v, Output: %s\n", err, output)
		panic("Failed to run the scaling tests")
	}
	fmt.Printf("Results of running tests: \n%s", output)
	fmt.Println("Finished running scale tests")
}

func runTiming(config TestConfig) {
	fmt.Println("Running Timing Tests")
	var testExec *exec.Cmd
	shouldPlot := fmt.Sprintf("-plot=%v", config.Plot)
	test := fmt.Sprintf("%v/benchmarktiming.test", config.Location)
	if config.Output != "" {
		testExec = exec.Command(test, shouldPlot)
	} else {
		output := fmt.Sprintf("-outFile=%v/%v", os.Getenv("PWD"), config.Output)
		testExec = exec.Command(test, output, shouldPlot)
	}
	output, err := testExec.Output()

	if err != nil && err.Error() != "exit status 197" {
		fmt.Printf("Error: %v, Output: %s", err, output)
		panic("Failed to run the timing tests")
	}
	fmt.Printf("Results of running tests: \n%s", output)
	fmt.Println("Finished timing scale tests")
}
