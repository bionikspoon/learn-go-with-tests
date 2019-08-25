package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type CountdownOperationsSpy struct {
	calls []string
}

func (spy *CountdownOperationsSpy) Sleep() {
	spy.calls = append(spy.calls, "sleep")
}
func (spy *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	spy.calls = append(spy.calls, "write")
	return
}

func TestCountdown(t *testing.T) {
	t.Run("3 2 1 Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		sleeper := &CountdownOperationsSpy{}

		Countdown(buffer, sleeper)

		got := buffer.String()
		want := `3
2
1
Go!
`

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("sleep before print", func(t *testing.T) {
		spy := &CountdownOperationsSpy{}

		Countdown(spy, spy)

		want := []string{"sleep", "write", "sleep", "write", "sleep", "write", "sleep", "write"}

		if !reflect.DeepEqual(spy.calls, want) {
			t.Errorf("wanted calls %v got %v", want, spy.calls)
		}
	})
}

type TimeSpy struct {
	durationSlept time.Duration
}

func (spy *TimeSpy) Sleep(duration time.Duration) {
	spy.durationSlept = duration
}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second
	timeSpy := &TimeSpy{}
	sleeper := ConfigurableSleeper{sleepTime, timeSpy.Sleep}

	sleeper.Sleep()

	if timeSpy.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, timeSpy.durationSlept)
	}
}
