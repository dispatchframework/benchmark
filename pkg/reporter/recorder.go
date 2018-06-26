package reporter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	. "github.com/logrusorgru/aurora"
)

type TimeRecord struct {
	Suite        string
	Records      map[string][]time.Duration
	Locks        map[string]*sync.RWMutex
	Measurements map[string]string
	Output       string
	shouldPlot   bool
	mu           *sync.RWMutex
	chart        func(map[string][]time.Duration, string)
}

func NewReporter(name string, output string, shouldPlot bool, charter func(map[string][]time.Duration, string)) *TimeRecord {
	var records TimeRecord
	var mu sync.RWMutex
	records.Suite = name
	records.Records = make(map[string][]time.Duration)
	records.Locks = make(map[string]*sync.RWMutex)
	records.mu = &mu
	records.Output = output
	records.shouldPlot = shouldPlot
	records.chart = charter
	return &records
}

func (t *TimeRecord) MakeMeasurement(name, value string) {
	t.Measurements[name] = value
}

func (t *TimeRecord) InitRecord(name string) {
	var lck sync.RWMutex
	var records []time.Duration
	t.Records[name] = records
	t.Locks[name] = &lck
}

func (t *TimeRecord) RecordTime(name string, length time.Duration) {
	lck := t.Locks[name]
	lck.Lock()
	defer lck.Unlock()
	records := t.Records[name]
	records = append(records, length)
	t.Records[name] = records
}

func (t *TimeRecord) GetRecord(name string) []time.Duration {
	lck := t.Locks[name]
	lck.Lock()
	defer lck.Unlock()
	return t.Records[name]
}

func GetStats(times []time.Duration) (float64, float64) {
	sum := 0.0
	for _, val := range times {
		sum += val.Seconds()
	}
	mean := sum / float64(len(times))

	distance := 0.0
	for _, val := range times {
		distance += math.Abs(mean - val.Seconds())
	}
	deviation := math.Sqrt(math.Pow(distance, 2) / float64(len(times)))
	return mean, deviation
}

func (t *TimeRecord) PrintResults() string {
	var result []string
	if len(t.Records) == 0 {
		return fmt.Sprintf(Sprintf("[%v]\n%v", t.Suite, Red("No Tests Run")))
	}
	t.OutToFile()
	for name, durations := range t.Records {
		sort.Slice(durations, func(i, j int) bool { return durations[i] > durations[j] })
		mean, sDev := GetStats(durations)
		field := Sprintf("Test: %v. %v Samples. \n\tSlowest: %v, \n\tFastest: %v. \n\tAverage: %v seconds. \n\tStandard Deviation: %v seconds.", name, len(durations), Red(durations[0]), Green(durations[len(durations)-1]), Cyan(mean), Magenta(Sprintf("%c %v", '\u00B1', sDev)))
		result = append(result, field)
	}
	result = append(result, "Individual Measurements")

	for field, value := range t.Measurements {
		result = append(result, Sprintf("%v: %s", Cyan(field), value))
	}

	if t.shouldPlot {
		name := fmt.Sprintf("chart-%v", t.Suite)
		fmt.Printf("Producing Chart %v\n", name)
		t.chart(t.Records, name)
	}
	return strings.Join(result, "\n")
}
