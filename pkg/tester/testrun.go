package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *Tester) RunSingleFunction(name string) {
	start := time.Now()
	util.ExecuteFunction(name)
	duration := time.Since(start)
	fmt.Printf("Single run: %v\n", duration.Seconds())
	t.aggregator.RecordValue("Run Single Function", duration.Seconds())
}

func (t *Tester) CheckFuncExists(name string) bool {
	for _, f := range t.functions {
		if f == name {
			return true
		}
	}
	return false
}

func (t *Tester) TestFuncRunSingle() {
	fmt.Println("Testing Run function")
	t.aggregator.InitRecord("Run Single Function")
	t.aggregator.AssignGraph("Execution", "Run Single Function")
	start := time.Now()
	function := "RunFuncTest"
	if !t.CheckFuncExists(function) {
		util.CreateFunction(function, testFunc)
		t.functions = append(t.functions, function)
	}
	for i := 0; i < samples; i++ {
		t.RunSingleFunction(function)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func (t *Tester) TestFuncRunSeries() {
	fmt.Println("Testing multiple function running in series")
	t.aggregator.InitRecord("Series Run Function")
	t.aggregator.AssignGraph("Execution", "Series Run Function")
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
		t.aggregator.RecordValue("Series Run Function", time.Since(start).Seconds())
	}
}

func (t *Tester) TestFuncRunParallel() {
	fmt.Println("Testing Multiple Function Execution in Parallel")
	function := "RunFuncTest"
	record := func(len float64) {
		t.aggregator.RecordValue("Parallel Run Function", len)
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
	t.aggregator.InitRecord("Parallel Run Function")
	t.aggregator.AssignGraph("Execution", "Parallel Run Function")
	for i := 0; i < samples; i++ {
		args := []string{function, testFunc}
		util.SyncRunRunners(toRun, record, 2, false, args...)
	}
}
