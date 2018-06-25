package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *TimeTests) TestRun(name string) {
	start := time.Now()
	util.ExecuteFunction(name, true)
	duration := time.Since(start)
	fmt.Printf("Single run: %v\n", duration.Seconds())
	aggregator.RecordTime("Run Single Function", duration)
}

func (t *TimeTests) TestFuncRunSingle() {
	fmt.Println("Testing Run function")
	aggregator.InitRecord("Run Single Function")
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
	aggregator.InitRecord("Series Run Function")
	if len(t.functions) <= 0 {
		util.CreateFunction("RunFuncTest", testFunc)
		t.functions = append(t.functions, "RunFuncTest")
	}
	start := time.Now()
	name := t.functions[0]
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			util.ExecuteFunction(name, true)
		}
		aggregator.RecordTime("Series Run Function", time.Since(start))
	}
}
