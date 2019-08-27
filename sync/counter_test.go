package sync

import (
	"sync"
	"testing"
)

func TestSync(t *testing.T) {
	t.Run("incrementing a counter", func(t *testing.T) {
		counter := NewCounter()

		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, 3)
	})

	t.Run("with concurrency", func(t *testing.T) {
		counter := NewCounter()

		want := 1000
		var wg sync.WaitGroup
		wg.Add(want)

		for i := 0; i < want; i++ {
			go func(w *sync.WaitGroup) {
				counter.Inc()
				w.Done()
			}(&wg)
		}

		wg.Wait()

		assertCounter(t, counter, want)
	})
}

func assertCounter(t *testing.T, counter *Counter, want int) {
	t.Helper()

	if got := counter.Value(); got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
