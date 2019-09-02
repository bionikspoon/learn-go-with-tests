package poker

import (
	"reflect"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	store := &StubPlayerStore{
		players: Players{},
	}

	t.Run("records Julia's win", func(t *testing.T) {
		in := strings.NewReader("Julia wins\n")
		cli := &CLI{store, in}

		cli.PlayPoker()

		want := []string{"RecordWin Julia"}
		if got := store.calledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("records Roger's win", func(t *testing.T) {
		in := strings.NewReader("Roger wins\n")
		cli := &CLI{store, in}

		cli.PlayPoker()

		want := []string{"RecordWin Julia", "RecordWin Roger"}
		if got := store.calledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})

}
