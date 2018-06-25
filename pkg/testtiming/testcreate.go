package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *TimeTests) MeasureSingleMake(name string) {
	fmt.Println("Creating Single Function")
	t.functions = append(t.functions, name)
	start := time.Now()
	util.CreateFunction(name, testFunc)
	aggregator.RecordTime("Make Function", time.Since(start))
}

func (t *TimeTests) TestFuncMakeSingle() {
	fmt.Println("Testing Make function")
	aggregator.InitRecord("Make Function")
	start := time.Now()
	for i := 0; i < samples; i++ {
		t.MeasureSingleMake(fmt.Sprintf("testFunc%v", i))
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func (t *TimeTests) TestFuncMakeSerial() {
	fmt.Println("Testing multiple function creation in series")
	aggregator.InitRecord("Series Function")
	start := time.Now()
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			name := fmt.Sprintf("Series%v-func%v", i, j)
			t.functions = append(t.functions, name)
			util.CreateFunction(name, testFunc)
		}
		aggregator.RecordTime("Series Function", time.Since(start))
	}
}

func (t *TimeTests) TestFuncMakeParallel() {
	fmt.Println("Testing multiple function creation in parallel")
	var funcs []string
	rcrd := func(len time.Duration) {
		aggregator.RecordTime("Parallel Function", len)
	}
	aggregator.InitRecord("Parallel Function")
	fmt.Printf("Samples: %v\n", samples)
	for i := 0; i < samples; i++ {
		funcs = util.SyncRunRunners(util.CreateFunction, testFunc, 2, rcrd, i)
		t.functions = append(t.functions, funcs...)
	}
}
