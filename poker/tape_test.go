package poker_test

import (
	"io/ioutil"
	"testing"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestTapeWrite(t *testing.T) {
	file, clean := poker.CreateTempFile(t, "12345")
	defer clean()

	tape := poker.NewTape(file)

	_, _ = tape.Write([]byte("abc"))

	_, _ = file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	want := "abc"
	if got := string(newFileContents); got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
