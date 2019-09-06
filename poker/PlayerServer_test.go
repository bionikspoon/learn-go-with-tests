package poker_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestPlayerServer(t *testing.T) {
	store := &poker.StubPlayerStore{
		Players: poker.Players{
			{1, "Joe Sixpack", 20},
			{2, "Jane User", 4},
			{3, "Creed", 3},
		},
	}
	server := poker.EnsurePlayerServer(t, store, &poker.SpyGame{})

	t.Run("show Joe Sixpack's score", func(t *testing.T) {
		request := poker.FetchShowScoreRequest("Joe Sixpack")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response, "20")

	})

	t.Run("show Jane User's score", func(t *testing.T) {
		request := poker.FetchShowScoreRequest("Jane User")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response, "4")
	})

	t.Run("show unknown users score", func(t *testing.T) {
		request := poker.FetchShowScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusNotFound)
		poker.AssertResponseBody(t, response, "0")
	})

	t.Run("it returns a 404 on missing players", func(t *testing.T) {
		request := poker.FetchShowScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusNotFound)
	})

	t.Run("it records wins on update", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.FetchUpdateScoreRequest("Apollo"))

		poker.AssertStatus(t, response, http.StatusAccepted)

		want := []string{
			"GetPlayerScore Joe Sixpack",
			"GetPlayerScore Jane User",
			"GetPlayerScore Floyd",
			"GetPlayerScore Apollo",
			"RecordWin Apollo",
		}

		if !reflect.DeepEqual(store.CalledWith, want) {
			t.Errorf("got %#v want %#v", store.CalledWith, want)
		}

	})

	t.Run("it returns OK on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.FetchIndexScoreRequest())

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertContentType(t, response, "application/json")
		poker.AssertLeague(t, server, store.Players)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := poker.EnsurePlayerServer(t, &poker.StubPlayerStore{}, &poker.SpyGame{})

		request := poker.FetchGameRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
	})

	t.Run("it prompts the user to enter number of players and starts the game", func(t *testing.T) {
		winner := "Ruth"

		game := &poker.SpyGame{}
		store := &poker.StubPlayerStore{}
		playerServer := poker.EnsurePlayerServer(t, store, game)
		server := httptest.NewServer(playerServer)

		defer server.Close()

		wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws := poker.EnsureWSDial(t, wsUrl)
		defer ws.Close()

		poker.EnsureWSWriteMessage(t, ws, "3")
		poker.EnsureWSWriteMessage(t, ws, winner)

		time.Sleep(20 * time.Millisecond)

		if got := game.StartedWith; got != 3 {
			t.Errorf("got %d want %d", got, 3)
		}

		if got := game.FinishedWith; got != winner {
			t.Errorf("got %q want %q", got, winner)
		}
	})
}
