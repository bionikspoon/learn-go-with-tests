package poker_test

import (
	"testing"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestRecordingWinsAndShowingThem(t *testing.T) {
	t.Run("InMemoryPlayerStore", func(t *testing.T) {

		server := poker.EnsurePlayerServer(t, poker.NewInMemoryPlayerStore(), &poker.SpyGame{})

		poker.AssertUpdateAndShow(t, server, "Pepper", 3)
		poker.AssertUpdateAndShow(t, server, "Candy", 6)
		poker.AssertUpdateAndShow(t, server, "Anne", 2)

		players := poker.Players{
			{0, "Candy", 6},
			{0, "Pepper", 3},
			{0, "Anne", 2},
		}
		poker.AssertLeague(t, server, players)
		poker.AssertLeague(t, server, players)

	})

	t.Run("DatabasePlayerStore", func(t *testing.T) {
		server := poker.EnsurePlayerServer(t, poker.NewDatabasePlayerStore(false, false), &poker.SpyGame{})

		poker.AssertUpdateAndShow(t, server, "Pepper", 3)
		poker.AssertUpdateAndShow(t, server, "Candy", 6)
		poker.AssertUpdateAndShow(t, server, "Anne", 2)

		players := poker.Players{
			{2, "Candy", 6},
			{1, "Pepper", 3},
			{3, "Anne", 2},
		}
		poker.AssertLeague(t, server, players)
		poker.AssertLeague(t, server, players)
	})

	t.Run("FileSystemPlayerStore", func(t *testing.T) {
		database, cleanup := poker.CreateTempFile(t, "")
		defer cleanup()

		server := poker.EnsurePlayerServer(t, poker.NewFileSystemPlayerStore(database), &poker.SpyGame{})

		poker.AssertUpdateAndShow(t, server, "Pepper", 3)
		poker.AssertUpdateAndShow(t, server, "Candy", 6)
		poker.AssertUpdateAndShow(t, server, "Anne", 2)

		players := poker.Players{
			{0, "Candy", 6},
			{0, "Pepper", 3},
			{0, "Anne", 2},
		}
		poker.AssertLeague(t, server, players)
		poker.AssertLeague(t, server, players)
	})

}
