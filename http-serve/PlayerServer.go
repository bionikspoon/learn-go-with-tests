package main

import (
	"fmt"
	"net/http"
)

type PlayerServer struct {
	store PlayerStore
}

func (server PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		server.create(w, player)
	case http.MethodGet:
		server.show(w, player)
	}

}

func (server PlayerServer) show(w http.ResponseWriter, player string) {

	score := server.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (server PlayerServer) create(w http.ResponseWriter, player string) {

	server.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
