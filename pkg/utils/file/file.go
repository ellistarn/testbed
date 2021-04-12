package file

import (
	"path/filepath"
	"runtime"
)

// Relative returns an absolute path relative to the provided path
func RelativeTo(path string) string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(file), path)
}
