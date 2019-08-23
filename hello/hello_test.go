package main

import "testing"

func TestHello(t *testing.T) {
	type args struct {
		name     string
		language Language
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"say hello to joe", args{"Joe Sixpack", English}, "Hello, Joe Sixpack!"},
		{"say hello to world", args{"World", English}, "Hello, World!"},
		{"say hello to world, given an empty string", args{"", English}, "Hello, World!"},
		{"say hello, in spanish", args{"Elodie", Spanish}, "Hola, Elodie!"},
		{"say hello, in french", args{"Camille", French}, "Bonjour, Camille!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hello(tt.args.name, tt.args.language); got != tt.want {
				t.Errorf("Hello() = %q, want %q", got, tt.want)
			}
		})
	}
}
