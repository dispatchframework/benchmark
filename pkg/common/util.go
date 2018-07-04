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

func CreateFunction(funcName, funcLocation string) {
	_, file := path.Split(funcLocation)
	handler := fmt.Sprintf("--handler=%s.handle", file[0:len(file)-len(path.Ext(funcLocation))])
	cmd := exec.Command("dispatch", "create", "function", funcName, funcLocation, "--image=python3", handler)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to create function %v. %v\n", funcName, err)
		panic("Unable to create function")
	}
	var fn struct {
		Status string
	}
	for fn.Status != "READY" {
		cmd := exec.Command("dispatch", "get", "functions", funcName, "--json")
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
	fmt.Printf("Created function %v\n", funcName)
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

func ExecuteFunction(funcName string) error {
	cmd := exec.Command("dispatch", "exec", funcName, "--wait")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Unable to run function: %v, \n%s\n", err, output)
		return err
	}
	var result ExecInfo
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Unable to unmarshal the result\n")
		return err
	}
	return nil
}

// func GetPodStatus(podName, ns string) string {
// 	var podStats struct {
// 		Status string
// 	}
// 	get := exec.Command("kubectl", "-n", ns, "describe", "pod", podName)
// 	output, err := get.CombinedOutput()
// 	if err != nil {
// 		log.Fatalf("Error getting status: %v. \n%s", err, output)
// 	}
// 	fmt.Printf("output: %s\n", output)
// 	json.Unmarshal(output, &podStats)
// 	return podStats.Status
// }

func SetupApi(name, target, path string) {
	fmt.Printf("Creating endpoint %v\n", name)
	createEndpoint := exec.Command("dispatch", "create", "api", "--method", "POST", "--path", fmt.Sprintf("/%v", path), name, target)
	output, err := createEndpoint.CombinedOutput()
	if err != nil {
		log.Fatalf("Unable to capture output. %v\n%v", err, output)
	}
	var status struct {
		Status string
	}
	for status.Status != "READY" {
		getStatus := exec.Command("dispatch", "get", "api", name, "--json")
		output, err := getStatus.Output()
		if err != nil {
			log.Fatalf("Unable to get status of the endpoint: %v\n", err)
		}
		if err := json.Unmarshal(output, &status); err != nil {
			log.Fatalf("Unable to decode the json status of endpoint, %v\n", err)
		}
	}
	fmt.Printf("Created Endpoint: %v\n", name)
}

func QueryApi(url, payload string) []byte {
	queryEndpoint := exec.Command(
		"curl", "-k", url,
		"-H", "Content-Type: application/json",
		"-d", payload,
	)
	output, err := queryEndpoint.CombinedOutput()
	if err != nil {
		log.Fatalf("Failure when querying api endpoint. %v\n%s", err, output)
	}
	return output
}
