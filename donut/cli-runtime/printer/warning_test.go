package printer_test

import (
	"bytes"
	"testing"

	"github.com/tomasbasham/donut/cli-runtime/printer"
)

func TestWarningPrinter_Print(t *testing.T) {
	tests := map[string]struct {
		msg      string
		expected string
		color    bool
	}{
		"message with color": {
			msg:      "this is a warning",
			expected: "\u001b[33;1mWarning:\u001b[0m this is a warning\n",
			color:    true,
		},
		"message without color": {
			msg:      "this is a warning",
			expected: "Warning: this is a warning\n",
			color:    false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			w := printer.NewWarningPrinter(buf, printer.WarningPrinterOptions{Color: tt.color})
			w.Print(tt.msg)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
