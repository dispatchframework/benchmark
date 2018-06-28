package reporter

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func (t *BenchmarkRecorder) ToJson(output *os.File) {
	type jsonIntermediate struct {
		Values      []float64 `json:"values"`
		Average     float64   `json:"average"`
		Measurement string    `json:"measurement"`
	}

	type BenchmarkReport struct {
		Tests     []jsonIntermediate `json:"tests"`
		Timestamp string             `json:"timestamp"`
	}
	var report BenchmarkReport
	report.Timestamp = fmt.Sprintf("%v", time.Now().Unix())
	for name, samples := range t.Records {
		var intermediate jsonIntermediate
		mean, _ := GetStats(samples)
		intermediate.Values = samples
		intermediate.Average = mean
		intermediate.Measurement = name
		report.Tests = append(report.Tests, intermediate)
	}
	body, err := json.Marshal(report)
	if err != nil {
		fmt.Printf("Unable to marshal json: %v\n", err)
		log.Fatal("Unable to jsonify")
	}
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	if _, err = writer.Write(body); err != nil {
		fmt.Printf("Unable to write json: %v\n", err)
		log.Fatal("Unable to write json")
	}

}

func (t *BenchmarkRecorder) ToCsv(output *os.File) {
	fmt.Printf("Outputting results to csv: %v\n", t.Output)
	writer := csv.NewWriter(output)
	defer writer.Flush()
	for name, samples := range t.Records {
		var data []string
		data = append(data, name)
		for _, val := range samples {
			data = append(data, fmt.Sprintf("%v", val))
		}
		mean, sDev := GetStats(samples)
		data = append(data, fmt.Sprintf("%v", mean))
		data = append(data, fmt.Sprintf("%v", sDev))
		if err := writer.Write(data); err != nil {
			fmt.Printf("Unable to write data to file: %v\n", err)
			log.Fatal("Unable to write data to file")
		}
	}

}

func (t *BenchmarkRecorder) OutToFile() {
	var output *os.File
	defer output.Close()
	if _, err := os.Stat(t.Output); err == nil {
		os.Remove(t.Output)
	}
	output, err := os.Create(t.Output)
	if err != nil {
		fmt.Printf("Unable to create new file, %v\n", err)
		log.Fatal("Unable to create new file")
	}
	switch ext := path.Ext(t.Output); ext {
	case ".json":
		t.ToJson(output)
	case ".csv":
		t.ToCsv(output)
	}

}
