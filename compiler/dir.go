package compiler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// A Dir implements FileSystem using the native file system restricted to a specific directory tree.
// While the FileSystem.Open method takes '/'-separated paths, a Dir's string value is a filename on the native file system, not a URL, so it is separated by filepath.Separator, which isn't necessarily '/'.
type Dir interface {
	Open(string) (io.Reader, error)
}

// Exposes a string as an unnamed file
type StringInputDir string

func (dir StringInputDir) Open(path string) (io.Reader, error) {
	if path != "" {
		return nil, fmt.Errorf("File not found: %s", path)
	}

	return strings.NewReader(string(dir)), nil
}

// Exposes operating system filesystem
type FsDir string

func (dir FsDir) Open(fp string) (io.Reader, error) {
	fp = filepath.Join(string(dir), fp)

	if !strings.HasPrefix(fp, "..") {
		return os.Open(fp)
	}

	return nil, fmt.Errorf("Invalid path: %s", fp)
}
