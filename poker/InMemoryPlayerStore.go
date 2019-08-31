package poker

import "sort"

type InMemoryPlayerStore struct {
	scores map[string]int
}

func (store *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *InMemoryPlayerStore) RecordWin(name string) {
	store.scores[name]++
}

func (store *InMemoryPlayerStore) GetLeague() (players Players) {
	for name, wins := range store.scores {
		players = append(players, Player{Name: name, Wins: wins})
	}
	sort.Sort(ByWins(players))
	return
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{make(map[string]int)}
}
