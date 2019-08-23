package main

import "fmt"

func main() {
	fmt.Println(Hello("World", English))
}

type Language int

const (
	Spanish Language = iota
	English
	French
)

const englishHello = "Hello"
const spanishHello = "Hola"
const frenchHello = "Bonjour"

func Hello(name string, language Language) string {
	if name == "" {
		name = "World"
	}

	return fmt.Sprintf("%v, %v!", prefix(language), name)
}

func prefix(language Language) string {
	switch language {
	case Spanish:
		return spanishHello

	case English:
		return englishHello

	case French:
		return frenchHello

	default:
		panic("unrecognized language")
	}
}
