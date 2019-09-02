package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	store, close, err := poker.NewFileSystemPlayerStoreFromFileName(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
