package workerpool_test

import (
	"sync"
	"testing"

	"testing/synctest"
	"time"

	"github.com/mufteev/workerpool"
)

// GOEXPERIMENT=synctest GOTRACEBACK=all go test

func TestPoolSimple(t *testing.T) {
	synctest.Run(func() {
		const (
			countWorkers   = 2
			countCollector = 2
		)

		wp, err := workerpool.NewPool(countWorkers, countCollector)
		if err != nil {
			t.Fatalf("workerpool: %s", err)
		}
		wp.RunBackground()

		timeout := 5 * time.Second
		now := time.Now()

		wg := sync.WaitGroup{}

		for range 2 {
			wg.Add(1)
			wp.AddTask(func() {
				defer wg.Done()
				time.Sleep(timeout)
			})
		}

		wg.Wait()
		wp.Stop()

		sub := time.Since(now)

		synctest.Wait()
		if sub > timeout {
			t.Fatalf("long: %d - %d", sub, timeout)
		}

		t.Logf("Ok: %d - %d", sub, timeout)
	})
}
