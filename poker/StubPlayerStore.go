package poker

import "fmt"

type StubPlayerStore struct {
	Players    Players
	CalledWith []string
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	store.CalledWith = append(store.CalledWith, fmt.Sprintf("GetPlayerScore %v", name))

	player := store.Players.Find(name)

	if player != nil {
		return player.Wins
	} else {
		return 0
	}
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.CalledWith = append(store.CalledWith, fmt.Sprintf("RecordWin %v", name))
}

func (store *StubPlayerStore) GetLeague() Players {
	store.CalledWith = append(store.CalledWith, "GetLeague")

	return store.Players
}
