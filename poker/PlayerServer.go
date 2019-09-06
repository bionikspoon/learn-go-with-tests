package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerServer struct {
	store    PlayerStore
	template *template.Template
	game     Game
	http.Handler
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	server := new(PlayerServer)

	templ, err := template.ParseFiles(RelativePath("game.html"))
	if err != nil {
		return nil, fmt.Errorf("problem loading template %v", err)
	}

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(server.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(server.playerHandler))
	router.Handle("/game", http.HandlerFunc(server.gameHandler))
	router.Handle("/ws", http.HandlerFunc(server.websocket))

	server.store = store
	server.template = templ
	server.Handler = router
	server.game = game
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

func (server *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	locals := struct{ Year int }{time.Now().Year()}

	if err := server.template.Execute(w, locals); err != nil {
		http.Error(w, fmt.Sprintf("problem rendering template %s", err.Error()), http.StatusInternalServerError)
	}
}

func (server *PlayerServer) websocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMessage := ws.WaitForMsg()
	numberOfPlayers, err := strconv.Atoi(numberOfPlayersMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not convert %q to number %v", numberOfPlayersMessage, err), http.StatusInternalServerError)
		return
	}
	server.game.Start(numberOfPlayers, ws)

	winner := ws.WaitForMsg()
	server.game.Finish(winner)

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

type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("semething went wrong %v", err), http.StatusInternalServerError)
		return nil
	}

	return &playerServerWS{ws}

}
func (ws *playerServerWS) WaitForMsg() string {
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(message)
}

func (ws *playerServerWS) Write(bytes []byte) (int, error) {
	if err := ws.WriteMessage(1, bytes); err != nil {
		return 0, err
	}
	return len(bytes), nil

}
