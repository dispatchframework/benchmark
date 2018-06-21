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

func (t *TimeRecord) ToJson(output *os.File) {
	type jsonIntermediate struct {
		Times []time.Duration
		Stats struct {
			Average float64
			StdDev  float64
		}
	}
	stats := make(map[string]jsonIntermediate)
	for name, durations := range t.Records {
		var intermediate jsonIntermediate
		mean, sDev := GetStats(durations)
		intermediate.Times = durations
		intermediate.Stats.Average = mean
		intermediate.Stats.StdDev = sDev
		stats[name] = intermediate
	}

	body, err := json.Marshal(stats)
	fmt.Println(stats["first"])
	fmt.Printf("%s\n", body)
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

func (t *TimeRecord) ToCsv(output *os.File) {
	fmt.Printf("Outputting results to csv: %v\n", t.Output)
	writer := csv.NewWriter(output)
	defer writer.Flush()
	for name, durations := range t.Records {
		var data []string
		data = append(data, name)
		for _, val := range durations {
			data = append(data, fmt.Sprintf("%v", val))
		}
		mean, sDev := GetStats(durations)
		data = append(data, fmt.Sprintf("%v", mean))
		data = append(data, fmt.Sprintf("%v", sDev))
		if err := writer.Write(data); err != nil {
			fmt.Printf("Unable to write data to file: %v\n", err)
			log.Fatal("Unable to write data to file")
		}
	}

}

func (t *TimeRecord) OutToFile() {
	var output *os.File
	defer output.Close()
	if _, err := os.Stat(t.Output); os.IsNotExist(err) {
		if output, err = os.Create(t.Output); err != nil {
			fmt.Printf("Unable to create new file, %v\n", err)
			log.Fatal("Unable to create new file")
		}
	} else {
		if output, err = os.Open(t.Output); err != nil {
			fmt.Printf("Unable to open file, %v\n", err)
			log.Fatal("Unable to open file")
		}
	}
	switch ext := path.Ext(t.Output); ext {
	case ".json":
		t.ToJson(output)
	case ".csv":
		t.ToCsv(output)
	}

}
