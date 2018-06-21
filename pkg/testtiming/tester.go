package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dispatchframework/benchmark/pkg/reporter"
)

var (
	aggregator *reporter.TimeRecord
)

type TimeTests int

func (t TimeTests) BenchmarkTest1() {
	aggregator.InitRecord("first")
	fmt.Println("Measuring first test")
	start := time.Now()
	for i := 0; i < 5; i++ {
		start = time.Now()
		fmt.Println("Hello World")
		aggregator.RecordTime("first", time.Since(start))
	}
}

func main() {
	aggregator = reporter.NewReporter("test runner")
	v := reflect.ValueOf(TimeTests(0))
	for k := 1; k < 2; k++ {
		v.MethodByName(fmt.Sprintf("BenchmarkTest%v", k)).Call(nil)
	}
	fmt.Printf("Values: %v\n", aggregator.GetRecord("first"))
	fmt.Println(aggregator.PrintResults())
}
