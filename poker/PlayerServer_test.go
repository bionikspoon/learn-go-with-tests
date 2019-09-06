package poker_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

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

		assertContentType(t, response, "application/json")
		poker.AssertStatus(t, response, http.StatusOK)
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
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &poker.SpyGame{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(poker.EnsurePlayerServer(t, &poker.StubPlayerStore{}, game))
		ws := ensureWSDial(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		ensureWSWriteMessage(t, ws, "3")
		ensureWSWriteMessage(t, ws, winner)

		assertGameStartedWith(t, game, 3)
		assertGameFinishedWith(t, game, winner)

		within(t, 10*time.Millisecond, func() {
			assertBlindAlert(t, ws, wantedBlindAlert)
		})

	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := response.Result().Header.Get("content-type"); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertGameStartedWith(t *testing.T, game *poker.SpyGame, want int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == want
	})

	if !passed {
		t.Errorf("got %q want %q", game.StartedWith, want)
	}
}

func assertGameFinishedWith(t *testing.T, game *poker.SpyGame, want string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == want
	})

	if !passed {
		t.Errorf("got %q want %q", game.FinishedWith, want)
	}
}

func assertBlindAlert(t *testing.T, ws *websocket.Conn, want string) {
	t.Helper()
	_, message, err := ws.ReadMessage()

	if err != nil {
		t.Error(err)
	}

	if got := string(message); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func within(t *testing.T, duration time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(duration):
		t.Error("timed out")
	case <-done:
	}
}

func retryUntil(duration time.Duration, fn func() bool) bool {
	deadline := time.Now().Add(duration)

	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
	}
	return false
}

func ensureWSDial(t *testing.T, url string) *websocket.Conn {
	t.Helper()

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open websocket server on %s %v", url, err)
	}

	return ws
}

func ensureWSWriteMessage(t *testing.T, ws *websocket.Conn, message string) {
	t.Helper()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
