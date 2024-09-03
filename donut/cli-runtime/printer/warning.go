package printer

import (
	"fmt"
	"io"
)

const (
	yellowColor = "\u001b[33;1m"
	resetColor  = "\u001b[0m"
)

type WarningPrinter struct {
	// out is the writer to output warnings to
	out io.Writer
	// opts contains options controlling warning output
	opts WarningPrinterOptions
}

// WarningPrinterOptions controls the behavior of a WarningPrinter constructed using NewWarningPrinter()
type WarningPrinterOptions struct {
	// Color indicates that warning output can include ANSI color codes
	Color bool
}

// NewWarningPrinter returns an implementation of warningPrinter that outputs warnings to the specified writer.
func NewWarningPrinter(out io.Writer, opts WarningPrinterOptions) *WarningPrinter {
	return &WarningPrinter{out, opts}
}

// Print prints warnings to the configured writer.
func (w *WarningPrinter) Print(message string) {
	if w.opts.Color {
		fmt.Fprintf(w.out, "%sWarning:%s %s\n", yellowColor, resetColor, message)
	} else {
		fmt.Fprintf(w.out, "Warning: %s\n", message)
	}
}
