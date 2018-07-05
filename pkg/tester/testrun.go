package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *Tester) runSingleFunction(name, measurement string) {
	start := time.Now()
	util.ExecuteFunction(name)
	duration := time.Since(start)
	fmt.Printf("Single run: %v\n", duration.Seconds())
	t.aggregator.RecordValue(measurement, duration.Seconds())
}

// CheckFuncExists checks whether or not a given function has already been created
func (t *Tester) CheckFuncExists(name string) bool {
	for _, f := range t.functions {
		if f == name {
			return true
		}
	}
	return false
}

// TestFuncRunSingle just runs a function once and measures how long it takes
func (t *Tester) TestFuncRunSingle() {
	fmt.Println("Testing Run function")
	measurement := "Run Single Function"
	t.aggregator.InitRecord("Run Single Function")
	t.aggregator.AssignGraph("Execution", measurement)
	start := time.Now()
	function := "RunFuncTest"
	if !t.CheckFuncExists(function) {
		util.CreateFunction(function, testFunc)
		t.functions = append(t.functions, function)
	}
	for i := 0; i < samples; i++ {
		t.runSingleFunction(function, measurement)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

// TestFuncRunSeries just runs a function multiple times in series and measures how long it takes
func (t *Tester) TestFuncRunSeries() {
	fmt.Println("Testing multiple function running in series")
	measurement := "Run Functions in Series"
	t.aggregator.InitRecord(measurement)
	t.aggregator.AssignGraph("Execution", measurement)
	function := "RunFuncTest"
	if !t.CheckFuncExists(function) {
		util.CreateFunction(function, testFunc)
		t.functions = append(t.functions, function)
	}
	start := time.Now()
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			util.ExecuteFunction(function)
		}
		t.aggregator.RecordValue(measurement, time.Since(start).Seconds())
	}
}

// TestFuncRunParallel just runs a function multiple times in parallel and measures how long it takes
func (t *Tester) TestFuncRunParallel() {
	fmt.Println("Testing Multiple Function Execution in Parallel")
	measurement := "Run Functions in Parallel"
	function := "RunFuncTest"
	record := func(len float64) {
		t.aggregator.RecordValue(measurement, len)
	}
	if !t.CheckFuncExists(function) {
		util.CreateFunction(function, testFunc)
		t.functions = append(t.functions, function)
	}
	toRun := func(args ...string) {
		if len(args) < 1 {
			panic("Not enough args to run function")
		}
		util.ExecuteFunction(function)
	}
	t.aggregator.InitRecord(measurement)
	t.aggregator.AssignGraph("Execution", measurement)
	for i := 0; i < samples; i++ {
		args := []string{function, testFunc}
		util.SyncRunRunners(toRun, record, 5, false, args...)
	}
}
