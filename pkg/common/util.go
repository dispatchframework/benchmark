package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"
)

// CreateFunction wraps dispatch's create function method, waits until function is ready to execute before it returns
func CreateFunction(funcName, funcLocation string) {
	_, file := path.Split(funcLocation)
	handler := fmt.Sprintf("--handler=%s.handle", file[0:len(file)-len(path.Ext(funcLocation))])
	cmd := exec.Command("dispatch", "create", "function", funcName, funcLocation, "--image=python3", handler)
  output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to create function %v. %v\n%s\n", funcName, err, output)
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
	fmt.Println("Running once to negate zero-scaling")
	if err := ExecuteFunction(funcName); err != nil {
		log.Fatalf("Failed to run function %v. %v", funcName, err)
	}
}

// DeleteFunction wraps dispatch's delete function command
func DeleteFunction(funcName string) error {
	cmd := exec.Command("dispatch", "delete", "function", funcName)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to delete function %v\n", funcName)
	}
	fmt.Printf("Delete function %v\n", funcName)
	return err
}

// ExecuteFunction runs a function once
func ExecuteFunction(funcName string) error {
	cmd := exec.Command("dispatch", "exec", funcName, "--wait")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Unable to run function: %v, \n%s\n", err, output)
		return err
	}
	return nil
}

// SetupAPI creates an api endpoint for the target function
func SetupAPI(name, target, path string) {
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

// QueryAPI curls the specified url, posting a payload
func QueryAPI(url, payload string) []byte {
	queryEndpoint := exec.Command(
		"curl", "-k", url,
		"-H", "Content-Type: application/json",
		"-d", payload,
	)
	output, err := queryEndpoint.CombinedOutput()
	if err != nil {
		log.Fatalf("Failure when querying API endpoint. %v\n%s", err, output)
	}
	return output
}
