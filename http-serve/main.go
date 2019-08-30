package main

import (
	"log"
	"net/http"
)

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

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Printf("Could not listen on port 5000 %v", err)
	}
}
