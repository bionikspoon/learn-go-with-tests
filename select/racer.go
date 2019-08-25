package racer

import (
	"fmt"
	"net/http"
	"time"
)

func Racer(a, b string) (winner string, err error) {
	return ConfigurableRacer(a, b, 10*time.Second)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, err error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("time out after waiting for %q and %q", a, b)
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)

	go func() {
		_, err := http.Get(url)
		if err != nil {
			ch <- false
		} else {
			ch <- true
		}
	}()

	return ch
}
