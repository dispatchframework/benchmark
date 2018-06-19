package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type TestConfig struct {
	ToRun    bool `json:"enabled"`
	Plot     bool
	Output   string `json:"omitempty"`
	Location string
	Samples  int
}

type RunnerConfig struct {
	Scaling TestConfig `json:"Scale,omitempty"`
	Timing  TestConfig `json:"Timing,omitempty"`
	Api     TestConfig `json:"Api,omitempty"`
}

func ReadJson(path string) RunnerConfig {
	fmt.Printf("Grabbing config from %s\n", path)
	rawConfig, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Unable to read provided configuration file")
	}
	var config RunnerConfig
	yaml.Unmarshal(rawConfig, &config)
	return config
}
