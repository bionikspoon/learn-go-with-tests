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
	"testing"
)

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

	got := decodePlayersFromResponse(t, response)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  %v\nwant %v", got, want)
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

func decodePlayersFromResponse(t *testing.T, response *httptest.ResponseRecorder) (players Players) {
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

func EnsurePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	t.Helper()

	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatalf("could not ensure player server %v", err)
	}

	return server
}
