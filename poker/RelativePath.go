package poker

import (
	"fmt"
	"path"
	"runtime"
)

func RelativePath(elem ...string) string {
	_, file, _, ok := runtime.Caller(1)

	if !ok {
		panic(fmt.Sprintf("something went wrong: %q", file))
	}

	elem = append([]string{path.Dir(file)}, elem...)

	return path.Join(elem...)
}
