package main

import (
	"fmt"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *Tester) MeasureSingleMake(name string) {
	fmt.Println("Creating Single Function")
	t.functions = append(t.functions, name)
	start := time.Now()
	util.CreateFunction(name, testFunc)
	t.aggregator.RecordValue("Make Function", time.Since(start).Seconds())
}

func (t *Tester) TestFuncMakeSingle() {
	fmt.Println("Testing Make function")
	t.aggregator.InitRecord("Make Function")
	start := time.Now()
	for i := 0; i < samples; i++ {
		t.MeasureSingleMake(fmt.Sprintf("testFunc%v", i))
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}

func (t *Tester) TestFuncMakeSerial() {
	fmt.Println("Testing multiple function creation in series")
	t.aggregator.InitRecord("Series Function")
	start := time.Now()
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			name := fmt.Sprintf("Series%v-func%v", i, j)
			t.functions = append(t.functions, name)
			util.CreateFunction(name, testFunc)
		}
		t.aggregator.RecordValue("Series Function", time.Since(start).Seconds())
	}
}

func (t *Tester) TestFuncMakeParallel() {
	fmt.Println("Testing Multiple Function Creation in Parallel")
	record := func(len float64) {
		t.aggregator.RecordValue("Parallel Function", len)
	}
	toRun := func(args ...string) {
		if len(args) < 2 {
			panic("Not enough args to create function")
		}
		name := args[0]
		location := args[1]
		t.functions = append(t.functions, name)
		util.CreateFunction(name, location)
	}
	t.aggregator.InitRecord("Parallel Function")
	for i := 0; i < samples; i++ {
		args := []string{fmt.Sprintf("parallel%v", i), testFunc}
		util.SyncRunRunners(toRun, record, 2, true, args...)
	}
}
