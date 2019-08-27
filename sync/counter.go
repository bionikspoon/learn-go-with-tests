package sync

import (
	"sync"
)

type Counter struct {
	value int
	mutex sync.Mutex
}

func (counter *Counter) Inc() {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	counter.value++

}

func (counter *Counter) Value() int {
	return counter.value
}

func NewCounter() *Counter {
	return &Counter{}
}
