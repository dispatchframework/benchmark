package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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

func (t *TimeTests) ApiBaseline() string {
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

func (t *TimeTests) CompareApiExec() {
	api := t.ApiBaseline()
	apis := exec.Command("dispatch", "get", "api")
	output, err := apis.CombinedOutput()
	if err != nil {
		log.Fatalf("Failure when getting apis. %v\n%v", err, output)
	}
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
		t.aggregator.RecordTime("Api Times", apiDuration)
		startExec := time.Now()
		_ = util.ExecuteFunction("TargetFunc")
		funcDuration := time.Since(startExec)
		t.aggregator.RecordTime("Function Times", funcDuration)
		t.aggregator.RecordTime("CompareApiExec",
			apiDuration-funcDuration)
	}
	fmt.Printf("Total time: %v", time.Since(start).Seconds())
}
