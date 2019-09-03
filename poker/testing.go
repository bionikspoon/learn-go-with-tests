package poker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type StubPlayerStore struct {
	Players    Players
	CalledWith []string
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	store.CalledWith = append(store.CalledWith, fmt.Sprintf("GetPlayerScore %v", name))

	player := store.Players.Find(name)

	if player != nil {
		return player.Wins
	} else {
		return 0
	}
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.CalledWith = append(store.CalledWith, fmt.Sprintf("RecordWin %v", name))
}

func (store *StubPlayerStore) GetLeague() Players {
	store.CalledWith = append(store.CalledWith, "GetLeague")

	return store.Players
}

func AssertUpdateAndShow(t *testing.T, server *PlayerServer, player string, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		server.ServeHTTP(httptest.NewRecorder(), FetchUpdateScoreRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, FetchShowScoreRequest(player))

	AssertStatus(t, response, http.StatusOK)
	AssertResponseBody(t, response, strconv.Itoa(count))
}

func AssertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()

	if got := response.Code; got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func AssertResponseBody(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := response.Body.String(); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func AssertLeague(t *testing.T, server *PlayerServer, want Players) {
	t.Helper()

	response := httptest.NewRecorder()
	server.ServeHTTP(response, FetchIndexScoreRequest())

	got := DecodePlayersFromResponse(t, response)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  %v\nwant %v", got, want)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := response.Result().Header.Get("content-type"); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func FetchIndexScoreRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/league", nil)
}

func FetchShowScoreRequest(name string) *http.Request {

	return httptest.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", url.PathEscape(name)), nil)
}

func FetchUpdateScoreRequest(name string) *http.Request {
	return httptest.NewRequest(http.MethodPut, fmt.Sprintf("/players/%s", url.PathEscape(name)), nil)
}

func FetchGameRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/game", nil)
}

func DecodePlayersFromResponse(t *testing.T, response *httptest.ResponseRecorder) (players Players) {
	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(&players); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}

	return
}

func CreateTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db.json")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	if _, err = tmpfile.Write([]byte(initialData)); err != nil {
		t.Fatalf("could not write initial data %v", err)
	}

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
