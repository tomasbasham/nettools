package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tomasbasham/donut/cli-runtime/flag"
	"github.com/tomasbasham/donut/cli-runtime/iooption"
	"github.com/tomasbasham/donut/cli-runtime/printer"

	"github.com/tomasbasham/donut/internal/cmd/flags"
)

type DonutOptions struct {
	Arguments []string
	IOStreams iooption.IOStreams
}

// NewRootCommand creates the `donut` command with default arguments.
func NewRootCommand() *cobra.Command {
	stream := iooption.NewDefaultIOStreams()

	return NewRootCommandWithArgs(DonutOptions{
		Arguments: os.Args,
		IOStreams: stream,
	})
}

// NewRootCommandWithArgs creates the `donut` command and its nested
// children.
func NewRootCommandWithArgs(o DonutOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "donut",
	}

	pflags := cmd.PersistentFlags()
	flags.RegisterPersistentFlags(pflags)

	printerOpts := printer.WarningPrinterOptions{Color: true}
	printer := printer.NewWarningPrinter(o.IOStreams.ErrOut, printerOpts)
	cmd.SetGlobalNormalizationFunc(flag.WarnWordSepNormalizeFunc(printer))

	// The globlal normalisation function ensures that all flags specified meet
	// the desired format, changing users' input if necessary.
	cmd.SetGlobalNormalizationFunc(flag.WordSepNormalizeFunc())

	cmd.AddCommand(NewLookupCommand())
	cmd.AddCommand(NewProxyCommand())

	return cmd
}
