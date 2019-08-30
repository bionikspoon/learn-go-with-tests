package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]Player{
			{1, "Joe Sixpack", 20},
			{2, "Jane User", 3},
			{3, "Creed", 3},
		},
		nil,
	}
	server := NewPlayerServer(store)

	t.Run("show Pepper's score", func(t *testing.T) {
		request := showScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")

	})

	t.Run("show Floyd's score", func(t *testing.T) {
		request := showScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "10")

	})

	t.Run("it returns a 404 on missing players", func(t *testing.T) {
		request := showScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})

	t.Run("it records wins on update", func(t *testing.T) {
		player := "Apollo"

		response := httptest.NewRecorder()
		server.ServeHTTP(response, updateScoreRequest(player))

		assertStatus(t, response, http.StatusAccepted)

		if want := []string{player}; !reflect.DeepEqual(store.calledWith, want) {
			t.Errorf("got %#v want %#v", store.calledWith, want)
		}

	})

	t.Run("it returns OK on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, indexScoreRequest())

		assertStatus(t, response, http.StatusOK)
		assertLeague(t, response, store.league)
		assertContentType(t, response, "application/json")

	})
}

func TestRecordingWinsAndShowingThem(t *testing.T) {
	t.Run("InMemoryPlayerStore", func(t *testing.T) {
		server := NewPlayerServer(NewInMemoryPlayerStore())

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, indexScoreRequest())
		assertLeague(t, response, []Player{
			{0, "Pepper", 3},
			{0, "Candy", 6},
			{0, "Anne", 2},
		})

	})

	t.Run("DatabasePlayerStore", func(t *testing.T) {
		server := NewPlayerServer(NewDatabasePlayerStore(true))

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, indexScoreRequest())
		assertLeague(t, response, []Player{
			{1, "Pepper", 3},
			{2, "Candy", 6},
			{3, "Anne", 2},
		})
	})
}

func assertUpdateAndShow(t *testing.T, server *PlayerServer, player string, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		server.ServeHTTP(httptest.NewRecorder(), updateScoreRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, showScoreRequest(player))

	assertStatus(t, response, http.StatusOK)
	assertResponseBody(t, response, strconv.Itoa(count))
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()

	if got := response.Code; got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := response.Body.String(); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertLeague(t *testing.T, response *httptest.ResponseRecorder, want []Player) {
	got := getLeagueFromResponse(t, response)

	sort.Sort(ByName(got))
	sort.Sort(ByName(want))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	if got := response.Result().Header.Get("content-type"); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func indexScoreRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}
func showScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func updateScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/players/%s", name), nil)
	return request
}
func getLeagueFromResponse(t *testing.T, response *httptest.ResponseRecorder) (league []Player) {

	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}

	return
}
