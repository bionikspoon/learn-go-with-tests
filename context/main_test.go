package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	data := "hello, world"
	t.Run("it returns data from the store", func(t *testing.T) {

		store := &SpyStore{response: data, t: t}
		server := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if got := response.Body.String(); got != data {
			t.Errorf("got %q want %q", got, data)
		}

	})

	t.Run("it tells the store to cancel if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		server := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingContext, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingContext)

		response := &SpyResponseWriter{}
		server.ServeHTTP(response, request)

		if response.written {
			t.Error("a response should not have been written")
		}
	})
}

type SpyResponseWriter struct {
	written bool
}

func (writer *SpyResponseWriter) Header() http.Header {
	writer.written = true
	return nil
}

func (writer *SpyResponseWriter) Write([]byte) (int, error) {
	writer.written = true
	return 0, nil
}

func (writer *SpyResponseWriter) WriteHeader(int) {
	writer.written = true

}

type SpyStore struct {
	response string
	t        *testing.T
}

func (store *SpyStore) Fetch(ctx context.Context) (string, error) {

	data := make(chan string, 1)

	go func() {
		var result string
		for _, char := range store.response {
			select {
			case <-ctx.Done():
				store.t.Log("spy store got cancelled")
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(char)
			}
		}

		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case d := <-data:
		return d, nil

	}
}
