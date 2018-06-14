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
	if config.Output != "" {
		testExec = exec.Command(test)
	} else {
		output := fmt.Sprintf("-outFile='./%v'", config.Output)
		testExec = exec.Command(test, output)
	}
	output, err := testExec.Output()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		panic("Failed to run the scaling tests")
	}
	fmt.Printf("Results of running tests: \n%s", output)
	fmt.Println("Finished running scale tests")
}

func runTiming(config TestConfig) {
	fmt.Println("Running Timing Tests")
	var testExec *exec.Cmd
	test := fmt.Sprintf("%v/benchmarktiming.test", config.Location)
	if config.Output != "" {
		testExec = exec.Command(test)
	} else {
		output := fmt.Sprintf("-outFile='./%v'", config.Output)
		testExec = exec.Command(test, output)
	}
	output, err := testExec.Output()
	if err != nil {
		panic("Failed to run the scaling tests")
	}
	fmt.Printf("Results of running tests: \n%s", output)
	fmt.Println("Finished timing scale tests")
}
