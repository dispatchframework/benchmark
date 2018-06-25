package common

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

// cmd := exec.Command("dispatch", "get", "function", name, "--json", "--insecure")
// fmt.Println("Making Request")
// output, err := cmd.Output()
// fmt.Println("Got response")
// if err != nil {
// 	panic("Couldn't get function info")
// }
// if err := json.Unmarshal(output, &fn); err != nil {
// 	panic("Unable to marshal json")
// }
// return fn.Status

func GetStatus(name string, done chan string, i int) {
	var fn struct {
		Status string
	}
	for fn.Status != "READY" {
		cmd := exec.Command("dispatch", "get", "functions", name, "--json")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := json.NewDecoder(stdout).Decode(&fn); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}
	done <- fn.Status
}

func CreateFunction(funcName, funcLocation string) {
	_, file := path.Split(funcLocation)
	handler := fmt.Sprintf("--handler=%s.handle", file[0:len(file)-len(path.Ext(funcLocation))])
	cmd := exec.Command("dispatch", "create", "function", funcName, funcLocation, "--image=python3", handler)
	fmt.Printf("Creating function %v\n", funcName)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to create function %v. %v\n", funcName, err)
		panic("Unable to create function")
	}

	i := 0
	done := make(chan string, 1)
	go GetStatus(funcName, done, i)
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("TIMEOUT")
		panic("REQUEST TIMEDOUT")
	case <-done:
		fmt.Printf("SUCCESS: CREATED FUNCTION %v\n", funcName)
	}
}

func DeleteFunction(funcName string) error {
	cmd := exec.Command("dispatch", "delete", "function", funcName)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to delete function %v\n", funcName)
	}
	fmt.Printf("Delete function %v\n", funcName)
	return err
}

func ExecuteFunction(funcName string, wait bool) error {
	shouldWait := ""
	if wait {
		shouldWait = "--wait"
	}
	fmt.Printf("Starting to run function: %v\n", funcName)
	cmd := exec.Command("dispatch", "exec", funcName, shouldWait)
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
	fmt.Println("Finished Running Function")
	return nil
}
