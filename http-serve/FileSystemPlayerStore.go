package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	players  Players
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	_, _ = file.Seek(0, 0)
	var players Players
	if err := json.NewDecoder(file).Decode(&players); err != nil {
		players = Players{}
	}
	database := json.NewEncoder(&tape{file})
	return &FileSystemPlayerStore{database, players}
}

func (store *FileSystemPlayerStore) GetLeague() Players {
	sort.Sort(ByWins(store.players))
	return store.players
}

func (store *FileSystemPlayerStore) GetPlayerScore(name string) int {
	if player := store.players.Find(name); player != nil {
		return player.Wins
	} else {
		return 0
	}
}

func (store *FileSystemPlayerStore) RecordWin(name string) {

	if player := store.players.Find(name); player != nil {
		player.Wins++
	} else {
		store.players = append(store.players, Player{Id: 0, Name: name, Wins: 1})
	}

	if err := store.database.Encode(store.players); err != nil {
		log.Printf("could not encode file err: %#+v\n", err)
	}
}
