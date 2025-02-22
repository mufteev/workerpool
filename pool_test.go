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

		wp := workerpool.NewPool(countWorkers, countCollector)
		wp.RunBackground()

		timeout := 1 * time.Second
		now := time.Now()

		wg := sync.WaitGroup{}

		for i := range 2 {
			_ = i
			wg.Add(1)
			task := workerpool.NewTask(func() {
				defer wg.Done()
				time.Sleep(timeout)
			})

			wp.AddTask(task)
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
