package common

import (
	"fmt"
	"sync"
	"time"
)

func Runner(toRun func(string, string), name, location string, wg, start *sync.WaitGroup) {
	wg.Done()
	start.Wait()
	toRun(name, location)
	wg.Done()
}

func SyncRunRunners(toRun func(string, string), location string, runners int, record func(time.Duration), iteration int) []string {
	var wg sync.WaitGroup
	var startWg sync.WaitGroup
	var funcs []string
	start := time.Now()
	startWg.Add(runners)
	// Creation Step
	for i := 0; i < runners; i++ {
		wg.Add(1)
		name := fmt.Sprintf("Parallel%v-%v", iteration, i)
		funcs = append(funcs, name)
		go Runner(toRun, name, location, &wg, &startWg)
	}
	wg.Wait()
	wg.Add(runners)
	start = time.Now()
	startWg.Add(-1 * runners)
	wg.Wait()
	record(time.Since(start))
	return funcs
}
