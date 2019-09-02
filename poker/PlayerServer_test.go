package poker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestPlayerServer(t *testing.T) {
	store := &StubPlayerStore{
		players: Players{
			{1, "Joe Sixpack", 20},
			{2, "Jane User", 4},
			{3, "Creed", 3},
		},
	}
	server := NewPlayerServer(store)

	t.Run("show Joe Sixpack's score", func(t *testing.T) {
		request := showScoreRequest("Joe Sixpack")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")

	})

	t.Run("show Jane User's score", func(t *testing.T) {
		request := showScoreRequest("Jane User")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, "4")
	})

	t.Run("show unknown users score", func(t *testing.T) {
		request := showScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
		assertResponseBody(t, response, "0")
	})

	t.Run("it returns a 404 on missing players", func(t *testing.T) {
		request := showScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})

	t.Run("it records wins on update", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, updateScoreRequest("Apollo"))

		assertStatus(t, response, http.StatusAccepted)

		want := []string{
			"GetPlayerScore Joe Sixpack",
			"GetPlayerScore Jane User",
			"GetPlayerScore Floyd",
			"GetPlayerScore Apollo",
			"RecordWin Apollo",
		}

		if !reflect.DeepEqual(store.calledWith, want) {
			t.Errorf("got %#v want %#v", store.calledWith, want)
		}

	})

	t.Run("it returns OK on /league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, indexScoreRequest())

		assertStatus(t, response, http.StatusOK)
		assertContentType(t, response, "application/json")
		assertLeague(t, server, store.players)
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

func assertLeague(t *testing.T, server *PlayerServer, want Players) {
	t.Helper()

	response := httptest.NewRecorder()
	server.ServeHTTP(response, indexScoreRequest())

	got := getLeagueFromResponse(t, response)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  %v\nwant %v", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()

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

func getLeagueFromResponse(t *testing.T, response *httptest.ResponseRecorder) (league Players) {
	t.Helper()

	if err := json.NewDecoder(response.Body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}

	return
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
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
