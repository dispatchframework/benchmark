package main

import (
	"fmt"
	"path"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

// TestScaleTimes runs a function many times in parallel, each iteration adding more executions
func (s *Tester) TestScaleTimes() {
	limit := 32
	fmt.Println("Comparing Performance as More functions are fun in parallel")
	function := "LimitTester"
	toRun := func(args ...string) {
		util.ExecuteFunction(function)
	}
	path, _ := path.Split(testFunc)
	funcPath := fmt.Sprintf("%vtest.py", path)
	util.CreateFunction(function, funcPath)
	s.functions = append(s.functions, function)
	util.ExecuteFunction(function)
	for j := 2; j <= limit; j *= 2 {
		measurement := fmt.Sprintf("Time to run %v functions in parallel", j)
		s.aggregator.InitRecord(measurement)
		s.aggregator.AssignGraph("Scale", measurement)
		record := func(len float64) {
			s.aggregator.RecordValue(measurement, len)
		}
		for i := 0; i < samples; i++ {
			args := []string{}
			util.SyncRunRunners(toRun, record, j, false, args...)
		}
	}
}
