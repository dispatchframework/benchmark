package common

import (
	"fmt"
	"sync"
	"time"
)

// Runner runs a given function once. Multiple runners will start at the same time
func Runner(toRun func(...string), wg, start *sync.WaitGroup, args ...string) {
	wg.Done()
	start.Wait()
	toRun(args...)
	wg.Done()
}

// SyncRunRunners creates and starts a number of runners. It also records values after those runners have finished
func SyncRunRunners(toRun func(...string), record func(float64), runners int, create bool, args ...string) {
	var wg sync.WaitGroup
	var startWg sync.WaitGroup
	start := time.Now()
	startWg.Add(runners)
	for i := 0; i < runners; i++ {
		wg.Add(1)
		newArgs := make([]string, len(args))
		copy(newArgs, args)
		if create {
			name := args[0]
			newArgs[0] = fmt.Sprintf("%v-%v", name, i)
		}
		go Runner(toRun, &wg, &startWg, newArgs...)
	}
	wg.Wait()
	wg.Add(runners)
	start = time.Now()
	startWg.Add(-1 * runners)
	wg.Wait()
	record(time.Since(start).Seconds())
}
