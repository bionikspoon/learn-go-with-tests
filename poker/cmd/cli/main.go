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

	game := poker.NewCLI(store, os.Stdin, os.Stdout, poker.BlindAlerterFunc(poker.StdOutAlerter))
	game.PlayPoker()

}
