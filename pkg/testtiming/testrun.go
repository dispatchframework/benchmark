package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *TimeTests) TestRun(name string) {
	start := time.Now()
	util.ExecuteFunction(name)
	duration := time.Since(start)
	fmt.Printf("Single run: %v\n", duration.Seconds())
	t.aggregator.RecordTime("Run Single Function", duration)
}

func (t *TimeTests) TestFuncRunSingle() {
	fmt.Println("Testing Run function")
	t.aggregator.InitRecord("Run Single Function")
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

func (t *TimeTests) TestFuncRunSeries() {
	fmt.Println("Testing multiple function running in series")
	t.aggregator.InitRecord("Series Run Function")
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
		t.aggregator.RecordTime("Series Run Function", time.Since(start))
	}
}

func (t *TimeTests) TestFuncRunParallel() {
	fmt.Println("Testing Multiple Function Creation in Parallel")
	record := func(len time.Duration) {
		t.aggregator.RecordTime("Parallel Run Function", len)
	}
	if len(t.functions) <= 0 {
		util.CreateFunction("RunFuncTest", testFunc)
		t.functions = append(t.functions, "RunFuncTest")
	}

	toRun := func(args ...string) {
		if len(args) < 1 {
			panic("Not enough args to create function")
		}
		name := args[0]
		util.ExecuteFunction(name)
	}
	t.aggregator.InitRecord("Parallel Run Function")
	for i := 0; i < samples; i++ {
		args := []string{"RunFuncTest", testFunc}
		util.SyncRunRunners(toRun, record, 2, false, args...)
	}
}
