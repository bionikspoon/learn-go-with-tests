package main

import (
	"log"
	"net/http"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

const dbFileName = "game.db.json"

func main() {
	filepath := poker.RelativePath("../../", dbFileName)
	store, close, err := poker.NewFileSystemPlayerStoreFromFileName(filepath)
	if err != nil {
		log.Fatalf("could not create store %v", err)
		return
	}
	defer close()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatalf("could not create server %v", err)
		return
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
