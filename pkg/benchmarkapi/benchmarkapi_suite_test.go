package benchmarkapi_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/nickaashoek/benchmark/pkg/benchmarkplot"
	. "github.com/nickaashoek/benchmark/pkg/benchmarkreporter"
	. "github.com/nickaashoek/benchmark/pkg/benchmarktiming"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Endpoint struct {
	Name string
	IP   string
	Port int
	Path string
}

var (
	wg         sync.WaitGroup
	testFunc   string
	outputFile string
	endpoint   Endpoint
	shouldPlot bool
)

func init() {
	flag.StringVar(&outputFile, "outFile", fmt.Sprintf("./output-%v.csv", time.Now().Unix()),
		"Controls where the output of the tests are written")
	flag.StringVar(&testFunc, "function",
		fmt.Sprintf("%v/src/github.com/dispatchframework/benchmark/resources/functions/test.py", os.Getenv("GOPATH")),
		"What function to use to test")
	flag.BoolVar(&shouldPlot, "plot", true, "Should a plot be produced")
}

func TestBenchmarkapi(t *testing.T) {
	RegisterFailHandler(Fail)
	reporter := NewDispatchReporter(outputFile, shouldPlot, ApiThroughputPlot)
	reporters := []Reporter{reporter}
	RunSpecsWithDefaultAndCustomReporters(t, "Dispatch Suite", reporters)
}

var _ = BeforeSuite(func() {
	// Setup the API Target Function
	setupWorker := Worker{Me: 1, Function: "api-target"}
	err := setupWorker.CreateFunction(testFunc)
	if err != nil {
		log.Fatal("Unable to setup initial function")
	}
	// Setup the API Endpoint
	// dispatch create api --method POST --path /hello post-hello hello-py

	endpoint = Endpoint{Name: "test-api", IP: "", Port: 80, Path: "testApi"}
	createEndpoint := exec.Command("dispatch", "create", "api", "--method", "POST", "--path", "/testApi", "test-api", "api-target")
	_, err = createEndpoint.Output()
	if err != nil {
		log.Fatalf("Unable to create the api endpoint, %v\n", err)
	}
	var status struct {
		Status string
	}
	for status.Status != "READY" {
		getStatus := exec.Command("dispatch", "get", "api", endpoint.Name, "--json")
		output, err := getStatus.Output()
		if err != nil {
			log.Fatalf("Unable to get status of the endpoint: %v\n", err)
		}
		if err := json.Unmarshal(output, &status); err != nil {
			log.Fatalf("Unable to decode the json status of endpoint, %v\n", err)
		}
	}

	// Run function once to mitigate riff zero scaling
	runFunction := exec.Command("dispatch", "exec", "api-target", "--wait")
	err = runFunction.Run()
	if err != nil {
		log.Fatalf("Unable to run the target function: %v\n", err)
	}

	log.Println("Finished creating the api endpoint")
})

var _ = AfterSuite(func() {
	// Tear down the API Endpoint and function
	deleteEndpoint := exec.Command("dispatch", "delete", "api", endpoint.Name)
	err := deleteEndpoint.Run()
	if err != nil {
		log.Fatalf("Unable to delete the api endpoint, %v\n", err)
	}
	deleteFunction := exec.Command("dispatch", "delete", "function", "api-target")
	err = deleteFunction.Run()
	if err != nil {
		log.Fatalf("Unable to delete the target function, %v\n", err)
	}
	log.Println("Finished the teardown")

})

func ManyQuery(increment *int64, call *exec.Cmd, wakeup *sync.Cond, wg *sync.WaitGroup, stop *int64) {
	wg.Done()
	wakeup.L.Lock()
	wakeup.Wait()
	defer wakeup.L.Unlock()
	for {
		if shouldStop := atomic.LoadInt64(stop); shouldStop != 0 {
			return
		}
		cmd := *call
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Unable to execute the api endpoint: %v\n", err)
		}
		atomic.AddInt64(increment, 1)
	}
}

func Flood(counter *int64, queriers int) {
	var wg sync.WaitGroup
	var stop int64
	var wakeLock sync.RWMutex
	payload := struct{ Value int }{2}
	jsonPayload, _ := json.Marshal(payload)
	wakeCond := sync.NewCond(wakeLock.RLocker())
	queryEndpoint := exec.Command(
		"curl", "-k", "http://35.203.141.195:80/testApi",
		"-H", "Content-Type: application/json",
		"-d", fmt.Sprintf("%s", jsonPayload),
	)
	wg.Add(queriers)
	for i := 0; i < queriers; i++ {
		go ManyQuery(counter, queryEndpoint, wakeCond, &wg, &stop)
	}
	wg.Wait()
	wakeCond.Broadcast()
	time.Sleep(1 * time.Second)
	atomic.StoreInt64(&stop, 1)
}

var _ = FDescribe("How many calls can we get through in 1 seconds", func() {
	// Want to get some idea of what throughput looks like
	Measure("2 queriers", func(b Benchmarker) {
		var counter int64
		Flood(&counter, 2)
		b.RecordValue("2", float64(counter))
	}, 10)
	Measure("4 queriers", func(b Benchmarker) {
		var counter int64
		Flood(&counter, 4)
		b.RecordValue("4", float64(counter))
	}, 10)
	Measure("8 queriers", func(b Benchmarker) {
		var counter int64
		Flood(&counter, 8)
		b.RecordValue("8", float64(counter))
	}, 10)

})

var _ = Describe("How long does it take to get a response from an API endpoint?", func() {
	Measure("How much overhead is there?", func(b Benchmarker) {
		functionRuntime := b.Time("Time for function to run", func() {
			runFunction := exec.Command("dispatch", "exec", "api-target", "--wait")
			err := runFunction.Run()
			if err != nil {
				log.Fatalf("Unable to run the target function: %v\n", err)
			}
		})
		var payload struct {
			Value int64
		}
		payload.Value = rand.Int63()
		jsonPayload, _ := json.Marshal(payload)
		apiRuntime := b.Time("Time to get response from api endpoint", func() {
			queryEndpoint := exec.Command(
				"curl", "-k", "http://35.203.141.195:80/testApi",
				"-H", "Content-Type: application/json",
				"-d", fmt.Sprintf("%s", jsonPayload),
			)
			_, err := queryEndpoint.Output()
			if err != nil {
				log.Fatalf("Query to API endpoint failed: %v\n", err)
			}
		})
		b.RecordValue("Approximate overhead of using the API endpoint", apiRuntime.Seconds()-functionRuntime.Seconds())
	}, 10)
})
