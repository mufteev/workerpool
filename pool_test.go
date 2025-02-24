package workerpool_test

import (
	"strings"
	"sync"
	"testing"

	"testing/synctest"
	"time"

	"github.com/mufteev/workerpool"
)

// GOEXPERIMENT=synctest GOTRACEBACK=all go test -v
// go test -benchmem -run ^$ -bench ^BenchmarkPool

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
			_ = wp.AddTask(func() {
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
	})
}

func TestPoolAddTaskNil(t *testing.T) {
	wp, err := workerpool.NewPool(1, 1)
	if err != nil {
		t.Fatalf("workerpool: %s", err)
	}

	wp.RunBackground()

	if err := wp.AddTask(nil); err != workerpool.ErrTaskNil {
		t.Fatalf("error: %s", err)
	}

	wp.Stop()
}

func TestPoolStop(t *testing.T) {
	wp, err := workerpool.NewPool(1, 1)
	if err != nil {
		t.Fatalf("workerpool: %s", err)
	}

	wp.RunBackground()
	wp.Stop()

	if err := wp.AddTask(func() {}); err != workerpool.ErrPoolStopped {
		t.Fatalf("error: %s", err)
	}
}

func BenchmarkPool_2Workers_2Collector(b *testing.B) {
	wp, _ := workerpool.NewPool(2, 2)
	wp.RunBackground()

	wg := sync.WaitGroup{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		_ = wp.AddTask(func() {
			defer wg.Done()

			var sb strings.Builder
			for j := 0; j < b.N; j++ {
				sb.WriteString("Hello, World!")
			}

			_ = sb.String()
		})
	}

	wg.Wait()
	wp.Stop()
}

func BenchmarkPool_10Workers_2Collector(b *testing.B) {
	wp, _ := workerpool.NewPool(10, 2)
	wp.RunBackground()

	wg := sync.WaitGroup{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		_ = wp.AddTask(func() {
			defer wg.Done()

			var sb strings.Builder
			for j := 0; j < b.N; j++ {
				sb.WriteString("Hello, World!")
			}

			_ = sb.String()
		})
	}

	wg.Wait()
	wp.Stop()
}
func BenchmarkPool_2Workers_10Collector(b *testing.B) {
	wp, _ := workerpool.NewPool(2, 10)
	wp.RunBackground()

	wg := sync.WaitGroup{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		_ = wp.AddTask(func() {
			defer wg.Done()

			var sb strings.Builder
			for j := 0; j < b.N; j++ {
				sb.WriteString("Hello, World!")
			}

			_ = sb.String()
		})
	}

	wg.Wait()
	wp.Stop()
}
