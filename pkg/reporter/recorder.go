package reporter

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	. "github.com/logrusorgru/aurora"
)

type TimeRecord struct {
	Suite   string
	Records map[string][]time.Duration
	Locks   map[string]*sync.RWMutex
	mu      *sync.RWMutex
}

func NewReporter(name string) *TimeRecord {
	var records TimeRecord
	var mu sync.RWMutex
	records.Suite = name
	records.Records = make(map[string][]time.Duration)
	records.Locks = make(map[string]*sync.RWMutex)
	records.mu = &mu
	return &records
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

func AverageRecords(times []time.Duration) float64 {
	sum := 0.0
	for _, val := range times {
		sum += val.Seconds()
	}
	return sum / float64(len(times))
}

func StdDevRecords(times []time.Duration, mean float64) float64 {
	sum := 0.0
	for _, val := range times {
		sum += math.Abs(mean - val.Seconds())
	}
	return math.Sqrt(math.Pow(sum, 2) / float64(len(times)))
}

func (t *TimeRecord) PrintResults() string {
	var result []string
	for name, durations := range t.Records {
		sort.Slice(durations, func(i, j int) bool { return durations[i] > durations[j] })
		mean := AverageRecords(durations)
		sDev := StdDevRecords(durations, mean)
		field := Sprintf("Test: %v. Slowest: %v, Fastest: %v. Average: %v seconds. Standard Deviation: %v seconds.", name, Red(durations[0]), Green(durations[len(durations)-1]), Cyan(mean), Magenta(Sprintf("%c %v", '\u00B1', sDev)))
		result = append(result, field)
	}
	return strings.Join(result, "\n")
}
