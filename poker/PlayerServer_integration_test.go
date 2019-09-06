package poker_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestRecordingWinsAndShowingThem(t *testing.T) {
	t.Run("InMemoryPlayerStore", func(t *testing.T) {

		server := poker.EnsurePlayerServer(t, poker.NewInMemoryPlayerStore(), &poker.SpyGame{})

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

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

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

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

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		players := poker.Players{
			{0, "Candy", 6},
			{0, "Pepper", 3},
			{0, "Anne", 2},
		}
		poker.AssertLeague(t, server, players)
		poker.AssertLeague(t, server, players)
	})

}

func assertUpdateAndShow(t *testing.T, server *poker.PlayerServer, player string, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		server.ServeHTTP(httptest.NewRecorder(), poker.FetchUpdateScoreRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, poker.FetchShowScoreRequest(player))

	poker.AssertStatus(t, response, http.StatusOK)
	poker.AssertResponseBody(t, response, strconv.Itoa(count))
}
