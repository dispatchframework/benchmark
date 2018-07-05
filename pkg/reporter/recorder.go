package reporter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
)

// BenchmarkRecorder maintains the records and mappings of a singe round of testing
// Multiple recorders may be used in one test suite or only one
// Locks are used to ensure parallel writing of records works
type BenchmarkRecorder struct {
	Suite        string
	Output       string
	mu           *sync.RWMutex
	Records      map[string][]float64
	Locks        map[string]*sync.RWMutex
	ChartCollect map[string][]string
	Graphs       map[string]func(map[string][]float64, string)
}

// NewReporter creates a new BenchmarkRecorder with the specified fields
func NewReporter(name string, output string) *BenchmarkRecorder {
	var records BenchmarkRecorder
	var mu sync.RWMutex
	records.Suite = name
	records.mu = &mu
	records.Output = output
	records.Records = make(map[string][]float64)
	records.Locks = make(map[string]*sync.RWMutex)
	records.ChartCollect = make(map[string][]string)
	records.Graphs = make(map[string]func(map[string][]float64, string))
	return &records
}

// InitRecord initializes a new record inside a recorder
func (t *BenchmarkRecorder) InitRecord(name string) {
	var lck sync.RWMutex
	var records []float64
	t.Records[name] = records
	t.Locks[name] = &lck
}

// RecordValue adds a value to a record
func (t *BenchmarkRecorder) RecordValue(name string, length float64) {
	lck := t.Locks[name]
	lck.Lock()
	defer lck.Unlock()
	records := t.Records[name]
	records = append(records, length)
	t.Records[name] = records
}

// AssignGraph assigns a record to a graph, which is useful when outputting multiple graphs
func (t *BenchmarkRecorder) AssignGraph(chart string, record string) {
	val, present := t.ChartCollect[chart]
	if present {
		val = append(val, record)
	} else {
		val = []string{record}
	}
	t.ChartCollect[chart] = val
}

// GetRecord returns the values in given record
func (t *BenchmarkRecorder) GetRecord(name string) []float64 {
	lck := t.Locks[name]
	lck.Lock()
	defer lck.Unlock()
	return t.Records[name]
}

// GetStats computes the average and stdDev of a record
func GetStats(records []float64) (float64, float64) {
	sum := 0.0
	for _, val := range records {
		sum += val
	}
	mean := sum / float64(len(records))

	distance := 0.0
	for _, val := range records {
		distance += math.Abs(mean - val)
	}
	deviation := math.Sqrt(math.Pow(distance, 2) / float64(len(records)))
	return mean, deviation
}

// PrintResults a nicely stringified representation of the results (with colors!)
func (t *BenchmarkRecorder) PrintResults() string {
	var result []string
	if len(t.Records) == 0 {
		return fmt.Sprintf(aurora.Sprintf("[%v]\n%v", t.Suite, aurora.Red("No Tests Run")))
	}
	t.OutToFile()
	for name, samples := range t.Records {
		sort.Slice(samples, func(i, j int) bool { return samples[i] > samples[j] })
		mean, sDev := GetStats(samples)
		field := aurora.Sprintf("Test: %v. %v Samples. \n\tSlowest: %v, \n\tFastest: %v. \n\tAverage: %v. \n\tStandard Deviation: %v",
			name, len(samples), aurora.Red(samples[0]), aurora.Green(samples[len(samples)-1]), aurora.Cyan(mean),
			aurora.Magenta(aurora.Sprintf("%c %v", '\u00B1', sDev)))
		result = append(result, field)
	}
	result = append(result, "Individual Measurements")
	for chart, records := range t.ChartCollect {
		grapher := t.Graphs[chart]
		recordMap := make(map[string][]float64)
		for _, record := range records {
			recordMap[record] = t.Records[record]
		}
		grapher(recordMap, chart)
	}
	return strings.Join(result, "\n")
}
