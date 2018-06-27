package common

import (
	"fmt"
	"sync"
	"time"
)

func Runner(toRun func(...string), wg, start *sync.WaitGroup, args ...string) {
	// func Runner(toRun func(...string), arg1, arg2 string, wg, start *sync.WaitGroup) {
	wg.Done()
	start.Wait()
	toRun(args...)
	wg.Done()
}

func SyncRunRunners(toRun func(...string), record func(float64), runners int, create bool, args ...string) {
	// func SyncRunRunners(toRun func(...string), location string, runners int, record func(time.Duration), iteration int) []string {
	var wg sync.WaitGroup
	var startWg sync.WaitGroup
	start := time.Now()
	startWg.Add(runners)
	// Creation Of Runners Step
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
