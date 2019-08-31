package main

import (
	"log"
	"net/http"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func main() {
	server := poker.NewPlayerServer(poker.NewInMemoryPlayerStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Printf("Could not listen on port 5000 %v", err)
	}
}
