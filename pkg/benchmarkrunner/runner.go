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
		runSuite(config.Timing, "Timing", "benchmarktiming")
	}
	if config.Scaling.ToRun {
		fmt.Println("Run scaling test")
		runSuite(config.Scaling, "Scale", "benchmarkscale")
	}
	if config.Api.ToRun {
		fmt.Println("Run api test")
		runSuite(config.Api, "Api", "benchmarkapi")
	}
}

func runSuite(config TestConfig, testName string, location string) {
	fmt.Printf("Running %v Tests: %v\n", testName, config)
	var testExec *exec.Cmd
	test := fmt.Sprintf("%v/%v.test", config.Location, location)
	shouldPlot := fmt.Sprintf("-plot=%v", config.Plot)
	samples := fmt.Sprintf("-samples=%v", config.Samples)
	if config.Output == "" {
		testExec = exec.Command(test, shouldPlot, samples)
	} else {
		output := fmt.Sprintf("-outFile=%v/%v", os.Getenv("PWD"), config.Output)
		testExec = exec.Command(test, output, shouldPlot, samples)
	}
	output, err := testExec.Output()

	// Exit status 197 is a special error status used by ginkgo to reflect programatic focus,
	// we don't want to report a test as failed in this case
	if err != nil && err.Error() != "exit status 197" {
		fmt.Printf("Error: %v, Output: %s\n", err, output)
		panic("Failed to run the test")
	}
	fmt.Printf("Results of running tests: \n%s", output)
	fmt.Printf("Finished running %v tests\n", testName)
}
