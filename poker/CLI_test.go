package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter number of players and starts the game", func(t *testing.T) {
		in := strings.NewReader("5\nJulia wins\n")
		out := &bytes.Buffer{}
		game := &SpyGame{}

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		if got := out.String(); got != poker.PlayerPrompt {
			t.Errorf("want %q got %q", got, poker.PlayerPrompt)
		}

		if got := game.StartedWith; got != 5 {
			t.Errorf("want %d got %d", got, 5)
		}

		if got := game.FinishedWith; got != "Julia" {
			t.Errorf("want %q got %q", got, "Julia")
		}
	})

	t.Run("given a non numerical value it prints an error and does not start the game", func(t *testing.T) {
		in := strings.NewReader("hello\n")
		out := &bytes.Buffer{}
		game := &SpyGame{}

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Error("game should not have started")
		}

		want := fmt.Sprintf("%s%s", poker.PlayerPrompt, "hello is not a number")
		if got := out.String(); got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

type SpyGame struct {
	StartedWith  int
	StartCalled  bool
	FinishedWith string
}

func (game *SpyGame) Start(startedWith int) {
	game.StartedWith = startedWith
	game.StartCalled = true
}
func (game *SpyGame) Finish(finishedWith string) {
	game.FinishedWith = finishedWith
}
