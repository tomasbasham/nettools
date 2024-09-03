package flag

import (
	"strings"

	"github.com/spf13/pflag"

	"github.com/tomasbasham/donut/cli-runtime/printer"
)

type normalizer func(f *pflag.FlagSet, name string) pflag.NormalizedName

type Printer interface {
	Print(string)
}

// wordSepNormalizeFunc changes all underscores in given CLI flags to
// hyphens. The StatusCake CLI will not accept underscores in flags, but instead
// of returning an error to the user we simply assume their intent and update
// their input.
func wordSepNormalizeFunc(p Printer) normalizer {
	return func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		if strings.Contains(name, "_") {
			p.Print("flag name " + name + " contains underscores, which are not allowed and have been replaced with hyphens")
			return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
		}
		return pflag.NormalizedName(name)
	}
}

// WordSepNormalizeFunc returns a normalizer that replaces underscores in flag
// names with hyphens.
func WordSepNormalizeFunc() normalizer {
	return wordSepNormalizeFunc(printer.Discard)
}

// WarnWordSepNormalizeFunc returns a normalizer that replaces underscores in
// flag names with hyphens and prints a warning message.
func WarnWordSepNormalizeFunc(p Printer) normalizer {
	return wordSepNormalizeFunc(p)
}
