package compiler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Dir interface {
	Open(string) (io.Reader, error)
}

type StringInputDir string

func (dir StringInputDir) Open(path string) (io.Reader, error) {
	if path != "" {
		return nil, fmt.Errorf("File not found: %s", path)
	}

	return strings.NewReader(string(dir)), nil
}

type FsDir string

func (dir FsDir) Open(fp string) (io.Reader, error) {
	fp = filepath.Join(string(dir), fp)

	if !strings.HasPrefix(fp, "..") {
		return os.Open(fp)
	}

	return nil, fmt.Errorf("Invalid path: %s", fp)
}
