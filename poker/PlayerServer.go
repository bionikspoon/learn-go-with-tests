package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerServer struct {
	store    PlayerStore
	template *template.Template
	http.Handler
}

func NewPlayerServer(store PlayerStore) (*PlayerServer, error) {
	server := new(PlayerServer)

	templ, err := template.ParseFiles(RelativePath("game.html"))
	if err != nil {
		return nil, fmt.Errorf("problem loading template %v", err)
	}

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(server.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(server.playerHandler))
	router.Handle("/game", http.HandlerFunc(server.game))
	router.Handle("/ws", http.HandlerFunc(server.websocket))

	server.store = store
	server.template = templ
	server.Handler = router
	return server, nil
}

func (server *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	league := server.store.GetLeague()
	if err := json.NewEncoder(w).Encode(league); err != nil {
		log.Printf("could not encode players err: %#+v\n", err)
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

func (server *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	locals := struct{ Year int }{time.Now().Year()}

	if err := server.template.Execute(w, locals); err != nil {
		http.Error(w, fmt.Sprintf("problem rendering template %s", err.Error()), http.StatusInternalServerError)
	}
}

func (server *PlayerServer) websocket(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("semething went wrong %v", err), http.StatusInternalServerError)
		return
	}

	_, message, err := ws.ReadMessage()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read ws message %v", err), http.StatusInternalServerError)
		return
	}

	server.store.RecordWin(string(message))

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
