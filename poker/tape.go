package poker

import "os"

type tape struct {
	file *os.File
}

func NewTape(file *os.File) *tape {
	return &tape{file}
}

func (tape tape) Write(bytes []byte) (n int, err error) {
	if err := tape.file.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := tape.file.Seek(0, 0); err != nil {
		return 0, err
	}

	return tape.file.Write(bytes)
}
