package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	tests := []struct {
		specify string
		name    string
		want    string
	}{
		{"it says hello to joe", "Joe Sixpack", "Hello, Joe Sixpack!"},
		{"it says hello to world", "World", "Hello, World!"},
	}
	for _, tt := range tests {
		t.Run(tt.specify, func(t *testing.T) {
			if got := Hello(tt.name); got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
