package main

import (
	"context"
	"fmt"
	"net/http"
)

func Server(store Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, err := store.Fetch(request.Context())

		if err != nil {
			return
		}

		fmt.Fprint(writer, data)
	}

}

type Store interface {
	Fetch(ctx context.Context) (string, error)
}
