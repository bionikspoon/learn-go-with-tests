package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestSHOWPlayer(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
	}
	server := &PlayerServer{store}

	t.Run("show Pepper's score", func(t *testing.T) {
		request := newShowScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")

	})

	t.Run("show Floyd's score", func(t *testing.T) {
		request := newShowScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "10")

	})

	t.Run("it returns a 404 on missing players", func(t *testing.T) {
		request := newShowScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Apollo"

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newPostScoreRequest(player))

		assertStatus(t, response, http.StatusAccepted)

		if want := []string{player}; !reflect.DeepEqual(store.calledWith, want) {
			t.Errorf("got %#v want %#v", store.calledWith, want)
		}

	})
}

func TestRecordingWinsAndShowingThem(t *testing.T) {
	t.Run("InMemoryPlayerStore", func(t *testing.T) {
		server := PlayerServer{NewInMemoryPlayerStore()}

		assertRecordAndShow(t, server, "Pepper", 3)
		assertRecordAndShow(t, server, "Candy", 6)
		assertRecordAndShow(t, server, "Anne", 2)

	})
	t.Run("DatabasePlayerStore", func(t *testing.T) {
		server := PlayerServer{NewDatabasePlayerStore(true)}

		assertRecordAndShow(t, server, "Pepper", 3)
		assertRecordAndShow(t, server, "Candy", 6)
		assertRecordAndShow(t, server, "Anne", 2)

	})
}

func assertRecordAndShow(t *testing.T, server PlayerServer, player string, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		server.ServeHTTP(httptest.NewRecorder(), newPostScoreRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newShowScoreRequest(player))

	assertStatus(t, response, http.StatusOK)
	assertResponseBody(t, response, strconv.Itoa(count))
}

func newShowScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newPostScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
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

type StubPlayerStore struct {
	scores     map[string]int
	calledWith []string
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.calledWith = append(store.calledWith, name)
}
