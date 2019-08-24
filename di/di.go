package main

import (
	"fmt"
	"io"
	"os"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %v!\n", name)

}

func main() {
	Greet(os.Stdout, "world")
}
