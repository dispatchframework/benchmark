package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *Tester) TestRun(name string) {
	start := time.Now()
	util.ExecuteFunction(name)
	duration := time.Since(start)
	fmt.Printf("Single run: %v\n", duration.Seconds())
	t.aggregator.RecordValue("Run Single Function", duration.Seconds())
}

func (t *Tester) TestFuncRunSingle() {
	fmt.Println("Testing Run function")
	t.aggregator.InitRecord("Run Single Function")
	t.aggregator.AssignGraph("Execution", "Run Single Function")
	start := time.Now()
	if len(t.functions) <= 0 {
		util.CreateFunction("RunFuncTest", testFunc)
		t.functions = append(t.functions, "RunFuncTest")
	}
	for i := 0; i < samples; i++ {
		name := t.functions[0]
		t.TestRun(name)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func (t *Tester) TestFuncRunSeries() {
	fmt.Println("Testing multiple function running in series")
	t.aggregator.InitRecord("Series Run Function")
	t.aggregator.AssignGraph("Execution", "Series Run Function")
	if len(t.functions) <= 0 {
		util.CreateFunction("RunFuncTest", testFunc)
		t.functions = append(t.functions, "RunFuncTest")
	}
	start := time.Now()
	name := t.functions[0]
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			util.ExecuteFunction(name)
		}
		t.aggregator.RecordValue("Series Run Function", time.Since(start).Seconds())
	}
}

func (t *Tester) TestFuncRunParallel() {
	fmt.Println("Testing Multiple Function Execution in Parallel")
	record := func(len float64) {
		t.aggregator.RecordValue("Parallel Run Function", len)
	}
	if len(t.functions) <= 0 {
		util.CreateFunction("RunFuncTest", testFunc)
		t.functions = append(t.functions, "RunFuncTest")
	}

	toRun := func(args ...string) {
		if len(args) < 1 {
			panic("Not enough args to run function")
		}
		name := args[0]
		util.ExecuteFunction(name)
	}
	t.aggregator.InitRecord("Parallel Run Function")
	t.aggregator.AssignGraph("Execution", "Parallel Run Function")
	for i := 0; i < samples; i++ {
		args := []string{"RunFuncTest", testFunc}
		util.SyncRunRunners(toRun, record, 2, false, args...)
	}
}
