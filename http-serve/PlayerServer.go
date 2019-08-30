package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Player struct {
	Id   int    `orm:"auto"`
	Name string `orm:"unique"`
	Wins int    `orm:"default(0)"`
}

type ByName []Player

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(server.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(server.playerHandler))

	server.Handler = router
	return server
}

func (server *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	league := server.store.GetLeague()
	if err := json.NewEncoder(w).Encode(league); err != nil {
		log.Printf("err: %#+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (server *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPut:
		server.update(w, player)
	case http.MethodGet:
		server.show(w, player)
	}
}

func (server *PlayerServer) show(w http.ResponseWriter, player string) {
	score := server.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (server *PlayerServer) update(w http.ResponseWriter, player string) {

	server.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
