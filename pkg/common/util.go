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

func CreateFunction(funcName, funcLocation string) error {
	_, file := path.Split(funcLocation)
	handler := fmt.Sprintf("--handler=%s.handle", file[0:len(file)-len(path.Ext(funcLocation))])
	cmd := exec.Command("dispatch", "create", "function", funcName, funcLocation, "--image=python3", handler)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to create function %v. %v\n", funcName, err)
		return err
	}
	var fn struct {
		Status string
	}
	for fn.Status != "READY" {
		cmd := exec.Command("dispatch", "get", "function", funcName, "--json")
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(output, &fn); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Created function %v\n", funcName)
	return nil
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
