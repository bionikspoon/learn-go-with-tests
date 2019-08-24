package di

import (
	"fmt"
	"io"
	"os"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %v!", name)

}

func Di() {
	Greet(os.Stdout, "world")
}
