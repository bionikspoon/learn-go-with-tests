package main

type InMemoryPlayerStore struct {
	scores map[string]int
}

func (store *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *InMemoryPlayerStore) RecordWin(name string) {
	store.scores[name]++
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{make(map[string]int)}
}
