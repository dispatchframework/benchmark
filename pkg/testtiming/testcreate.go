package main

import (
	"fmt"
	"log"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func (t *TimeTests) MeasureSingleMake(name string) {
	fmt.Println("Creating Single Function")
	t.functions = append(t.functions, name)
	start := time.Now()
	err := util.CreateFunction(name, testFunc)
	if err != nil {
		log.Println("Unable to create function!")
	}
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
	fmt.Println("Testing multiple functions in series")
	aggregator.InitRecord("Series Function")
	start := time.Now()
	for i := 0; i < samples; i++ {
		start = time.Now()
		for j := 0; j < 5; j++ {
			name := fmt.Sprintf("Series%v-func%v", i, j)
			t.functions = append(t.functions, name)
			err := util.CreateFunction(name, testFunc)
			if err != nil {
				fmt.Printf("Failed to make function %v. Error: %v\n", name, err)
				panic("Unable to create function")
			}
		}
		aggregator.RecordTime("Series Function", time.Since(start))
	}
}
