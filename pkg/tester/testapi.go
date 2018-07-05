package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	util "github.com/dispatchframework/benchmark/pkg/common"
)

func init() {
	// Get the ip for API
	fmt.Println("[Getting the IP address of the API server]")
	query := exec.Command("kubectl", "-n", "kong",
		"get", "service", "api-gateway-kongproxy", "-o=json")
	output, err := query.CombinedOutput()
	if err != nil {
		log.Fatalf("Failure when getting ip. %v\n%v", err, output)
	}
	var IPstruct struct {
		Status struct {
			LoadBalancer struct {
				Ingress []struct {
          IP string
				}
			}
		}
	}
	err = json.Unmarshal(output, &IPstruct)
	if err != nil {
		log.Fatalf("Unable to unmarshal json %v\n", err)
	}
	apiIP = IPstruct.Status.LoadBalancer.Ingress[0].IP
	fmt.Printf("Got the IP: %v\n", apiIP)
}

func (t *Tester) setupAPIBaseline() string {
	api := "testAPI"
	if !t.CheckFuncExists("TargetFunc") {
		util.CreateFunction("TargetFunc", testFunc)
		t.functions = append(t.functions, "TargetFunc")
		t.apis = append(t.apis, api)
		util.SetupAPI(api, "TargetFunc", api)
	}
	return api
}

// TestAPIvsExec measures the difference between doing dispatch exec vs. curl-ing an api endpoint
func (t *Tester) TestAPIvsExec() {
	api := t.setupAPIBaseline()
	measurement := "Measuring API Runtimes vs. Direct function execution times"
	t.aggregator.InitRecord(measurement)
	t.aggregator.InitRecord("API Execution Times")
	t.aggregator.InitRecord("Function Execution Times")
	start := time.Now()
	url := fmt.Sprintf("http://%v:80/%v", apiIP, api)
	for i := 0; i < samples; i++ {
		startAPI := time.Now()
		_ = util.QueryAPI(url, `{}`)
		apiDuration := time.Since(startAPI)
		t.aggregator.RecordValue("API Execution Times", apiDuration.Seconds())
		startExec := time.Now()
		_ = util.ExecuteFunction("TargetFunc")
		funcDuration := time.Since(startExec)
		t.aggregator.RecordValue("Function Execution Times", funcDuration.Seconds())
		t.aggregator.RecordValue(measurement,
			apiDuration.Seconds()-funcDuration.Seconds())
	}
	fmt.Printf("Total time: %v", time.Since(start).Seconds())
}

// TestAPIThroughput measures how many function requests can make it through in 1 second
func (t *Tester) TestAPIThroughput() {
	var wg sync.WaitGroup
	maxRunners := 6
	api := t.setupAPIBaseline()
	record := func(value float64) {
		return
	}
	counter := 0
	stop := int64(0)
	url := fmt.Sprintf("http://%v:80/%v", apiIP, api)
	_ = util.QueryAPI(url, `{}`)
	toRun := func(args ...string) {
		wg.Done()
		for {
			val := atomic.LoadInt64(&stop)
			if val != 0 {
				break
			}
			util.QueryAPI(url, `{}`)
			counter++
		}
	}
	for j := 1; j < maxRunners; j++ {
		runners := j
		test := fmt.Sprintf("Number of Functions run in 1 Second with %v parallel runners", runners)
		t.aggregator.InitRecord(test)
		t.aggregator.AssignGraph("Api", test)
		stop = 0
		counter = 0
		for i := 0; i < samples; i++ {
			args := []string{}
			wg.Add(runners)
			go util.SyncRunRunners(toRun, record, runners, false, args...)
			wg.Wait()
			select {
			case <-time.After(1 * time.Second):
				atomic.StoreInt64(&stop, 1)
				len, _ := time.ParseDuration(fmt.Sprintf("%vs", counter))
				t.aggregator.RecordValue(test, len.Seconds())
			}
		}
	}
}
