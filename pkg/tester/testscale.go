package main

import (
	"fmt"
	"path"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

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
		measurement := fmt.Sprintf("Scale Test %v runners", j)
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

// This may be something to come back to, I'm hitting limits on my computer before evicting the pod
// func (s *ScaleTests) GetLimit() {
// 	fmt.Println("Testing Multiple Function Creation in Parallel")
// 	toRun := func(args ...string) {
// 		util.ExecuteFunction("LimitTester")
// 	}
// 	if len(s.functions) <= 0 {
// 		function := fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH"))
// 		util.CreateFunction("LimitTester", function)
// 		s.functions = append(s.functions, "LimitTester")
// 		util.ExecuteFunction("LimitTester")
// 	}
// 	// First step is figuring out which pod we actually care about
// 	findPod := exec.Command("kubectl", "get", "pods", "-n", "dispatch")
// 	output, err := findPod.CombinedOutput()
// 	if err != nil {
// 		log.Fatalf("Failed to get pods, %v\n%s", err, output)
// 	}
// 	rx, err := regexp.Compile(`\S+function\-manager\S+`)
// 	pod := fmt.Sprintf("%s", rx.Find(output))
// 	fmt.Println(util.GetPodStatus(pod, "dispatch"))
// }
