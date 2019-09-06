package main

import (
	"fmt"
	"log"
	"os"

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

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
