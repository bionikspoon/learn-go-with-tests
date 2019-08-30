package main

type StubPlayerStore struct {
	scores     map[string]int
	calledWith []string
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.calledWith = append(store.calledWith, name)
}
