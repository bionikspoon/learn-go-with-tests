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
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err = http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
