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
					Ip string
				}
			}
		}
	}
	err = json.Unmarshal(output, &IPstruct)
	if err != nil {
		log.Fatalf("Unable to unmarshal json %v\n", err)
	}
	apiIP = IPstruct.Status.LoadBalancer.Ingress[0].Ip
	fmt.Printf("Got the IP: %v\n", apiIP)
}

func (t *Tester) ApiBaseline() string {
	var api string
	if len(t.functions) <= 0 {
		util.CreateFunction("TargetFunc", testFunc)
		t.functions = append(t.functions, "TargetFunc")
	}

	if len(t.apis) > 0 {
		api = t.apis[0]
	} else {
		api = "testApi"
		t.apis = append(t.apis, api)
		util.SetupApi(api, "TargetFunc", api)
	}
	return api
}

func (t *Tester) CompareApiExec() {
	api := t.ApiBaseline()
	t.aggregator.InitRecord("CompareApiExec")
	t.aggregator.InitRecord("Api Times")
	t.aggregator.InitRecord("Function Times")
	start := time.Now()
	url := fmt.Sprintf("http://%v:80/%v", apiIP, api)
	_ = util.ExecuteFunction("TargetFunc")
	for i := 0; i < samples; i++ {
		startApi := time.Now()
		_ = util.QueryApi(url, `{}`)
		apiDuration := time.Since(startApi)
		t.aggregator.RecordValue("Api Times", apiDuration.Seconds())
		startExec := time.Now()
		_ = util.ExecuteFunction("TargetFunc")
		funcDuration := time.Since(startExec)
		t.aggregator.RecordValue("Function Times", funcDuration.Seconds())
		t.aggregator.RecordValue("CompareApiExec",
			apiDuration.Seconds()-funcDuration.Seconds())
	}
	fmt.Printf("Total time: %v", time.Since(start).Seconds())
}

func (t *Tester) ApiMeasureThroughput() {
	var wg sync.WaitGroup
	maxRunners := 6
	api := t.ApiBaseline()
	record := func(value float64) {
		return
	}
	counter := 0
	stop := int64(0)
	url := fmt.Sprintf("http://%v:80/%v", apiIP, api)
	_ = util.QueryApi(url, `{}`)
	toRun := func(args ...string) {
		wg.Done()
		for {
			val := atomic.LoadInt64(&stop)
			if val != 0 {
				break
			}
			util.QueryApi(url, `{}`)
			counter++
		}
	}
	for j := 1; j < maxRunners; j++ {
		runners := j
		test := fmt.Sprintf("Functions in 1 Second with %v Runners", runners)
		t.aggregator.InitRecord(test)
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
