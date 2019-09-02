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

	database, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("could not open file err: %#+v\n", err)
	}

	store := poker.NewFileSystemPlayerStore(database)

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()

}
