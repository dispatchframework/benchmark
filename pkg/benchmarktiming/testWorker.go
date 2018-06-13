package benchmarktiming

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"path"
	"time"
)

type ExecInfo struct {
	ExecutedTime float64
	FinishedTime float64
	Name         string
}

type Run struct {
	Name         string
	ExecutedTime float64
	FinishedTime float64
	Status       string
	Logs         struct {
		Stderr string
		Stdout []string
	}
}

type Worker struct {
	Me       int
	Names    []string
	Function string
}

func RandomName(n int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

func (wk *Worker) CreateFunction(funcLocation string) error {
	_, file := path.Split(funcLocation)
	handler := fmt.Sprintf("--handler=%s.handle", file[0:len(file)-len(path.Ext(funcLocation))])
	cmd := exec.Command("dispatch", "create", "function", wk.Function, funcLocation, "--image=python3", handler)
	_, err := cmd.Output()
	if err != nil {
		DPrintf("Error in creating function: %v", err)
	} else {
		DPrintf("Worker %v finished creating function %v", wk.Me, wk.Function)
	}
	var fn struct {
		Status string
	}
	for fn.Status != "READY" {
		cmd := exec.Command("dispatch", "get", "function", wk.Function, "--json")
		output, err := cmd.Output()
		if err != nil {
			DPrintf("Uh Oh, failed to get information on fn\n")
			return err
		}
		if err := json.Unmarshal(output, &fn); err != nil {
			DPrintf("Failed to unmarshal function\n")
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	DPrintf("Function is now ready\n")
	return nil
}

func (wk *Worker) DeleteFunction() error {
	DPrintf("Deleting function\n")
	cmd := exec.Command("dispatch", "delete", "function", wk.Function)
	_, err := cmd.Output()
	if err != nil {
		DPrintf("Error in deleting function\n")
	}
	DPrintf("Deleted Function\n")
	return err
}

func (wk *Worker) ExecuteFunction(wait bool) error {
	// cmd := exec.Command("./createFunc.sh", "script-func", "./functions/test.py", "--image=python3", "--handler=test.handle")
	shouldWait := ""
	if wait {
		shouldWait = "--wait"
	}
	DPrintf("Worker %v running function: %v", wk.Me, wk.Function)
	cmd := exec.Command("dispatch", "exec", wk.Function, shouldWait)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to run function: %v, %v\n", err, output)
		return err
	}
	var result ExecInfo
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Unable to unmarshal the result\n")
		return err
	}
	return nil
}

func (wk *Worker) RecollectInfo() {
	for _, name := range wk.Names {
		var run Run
		DPrintf("Worker %v: Collecting information about run %v", wk.Me, name)
		for run.Status != "READY" {
			cmd := exec.Command("dispatch", "get", "runs", wk.Function, name, "--json")
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Uh Oh, failed to get information on run")
				log.Fatal(err)
			}
			if err := json.Unmarshal(output, &run); err != nil {
				fmt.Println("Failed to unmarshal")
				log.Fatal(err)
			}
			fmt.Printf("Worker %v: Collected info about %v, status: %v\n", wk.Me, name, run.Status)
		}
		fmt.Printf("Run %v took %v seconds\n", run.Name, run.FinishedTime-run.ExecutedTime)
	}
	cmd := exec.Command("dispatch", "delete", "function", wk.Function)
	_, err := cmd.Output()
	if err != nil {
		DPrintf("Worker %v unable to tear down function %v!", wk.Me, wk.Function)
	} else {
		DPrintf("Worker %v tore down function %v!", wk.Me, wk.Function)
	}
}
