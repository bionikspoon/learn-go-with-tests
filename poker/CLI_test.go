package poker_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestCLI(t *testing.T) {
	store := &poker.StubPlayerStore{
		Players: poker.Players{},
	}

	t.Run("records Julia's win", func(t *testing.T) {
		in := strings.NewReader("Julia wins\n")
		cli := poker.NewCLI(store, in)

		cli.PlayPoker()

		want := []string{"RecordWin Julia"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("records Roger's win", func(t *testing.T) {
		in := strings.NewReader("Roger wins\n")
		cli := poker.NewCLI(store, in)

		cli.PlayPoker()

		want := []string{"RecordWin Julia", "RecordWin Roger"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})

}
