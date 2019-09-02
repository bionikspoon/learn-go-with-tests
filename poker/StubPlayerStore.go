package poker

import (
	"fmt"
)

type StubPlayerStore struct {
	players    Players
	calledWith []string
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	store.calledWith = append(store.calledWith, fmt.Sprintf("GetPlayerScore %v", name))

	player := store.players.Find(name)

	if player != nil {
		return player.Wins
	} else {
		return 0
	}
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.calledWith = append(store.calledWith, fmt.Sprintf("RecordWin %v", name))
}

func (store *StubPlayerStore) GetLeague() Players {
	store.calledWith = append(store.calledWith, "GetLeague")

	return store.players
}
