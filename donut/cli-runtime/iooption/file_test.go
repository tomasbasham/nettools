package iooption_test

import (
	"testing"

	io "github.com/tomasbasham/donut/cli-runtime/iooption"
)

func TestOpenFile(t *testing.T) {
	// create a table based test
	tests := map[string]struct {
		filename string
		err      bool
	}{
		"open file": {
			filename: "testdata/file.txt",
			err:      false,
		},
		"open stdin": {
			filename: "-",
			err:      false,
		},
		"open missing file": {
			filename: "testdata/missing.txt",
			err:      true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := io.OpenFile(tt.filename)
			if tt.err && err == nil {
				t.Fatal("expected an error")
			}
			if !tt.err && err != nil {
				t.Fatal(err)
			}
		})
	}
}
