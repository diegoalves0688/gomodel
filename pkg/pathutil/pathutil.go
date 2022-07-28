package pathutil

import (
	"path/filepath"
	"runtime"
)

func RelativePath() string {
	// by passing "1" to "runtime.Caller(1)", the returned
	// valued is the relative path to the file containing
	// the call to "func RelativePath() string".
	_, b, _, _ := runtime.Caller(1) // nolint:dogsled
	return filepath.Dir(b)
}
