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

// ToJSON outputs all the recorded values to a JSON file. Intended use is for uploading to GCloud for parsing
func (t *BenchmarkRecorder) ToJSON(output *os.File) {
	type JSONIntermediate struct {
		Values      []float64 `json:"values"`
		Average     float64   `json:"average"`
		Measurement string    `json:"measurement"`
	}

	type BenchmarkReport struct {
		Tests     []JSONIntermediate `json:"tests"`
		Timestamp string             `json:"timestamp"`
	}
	var report BenchmarkReport
	report.Timestamp = fmt.Sprintf("%v", time.Now().Unix())
	for name, samples := range t.Records {
		var intermediate JSONIntermediate
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

// ToCsv outputs the recorded values as a csv file. Intended use is for later graphing/data analysis
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

// OutToFile outputs the results to a file, either csv or JSON
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
		t.ToJSON(output)
	case ".csv":
		t.ToCsv(output)
	}

}
