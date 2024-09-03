package iooption

import (
	"io"
	"os"
)

func OpenFile(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}
