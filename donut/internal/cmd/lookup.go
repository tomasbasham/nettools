package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tomasbasham/donut"
)

func NewLookupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "lookup",
		Short: "Lookup a domain name",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			resolver := donut.New(donut.GoogleHost)

			fqdn := args[0]

			t := donut.A
			if len(args) == 2 {
				t = getType(args[1])
			}

			question := donut.Question{
				FQDN:  fqdn,
				Type:  t,
				Class: donut.IN,
			}

			answer, err := resolver.Lookup(question)
			if err != nil {
				panic(err)
			}

			b, err := json.MarshalIndent(answer, "", "  ")
			if err != nil {
				panic(err)
			}

			fmt.Println(string(b))
		},
	}
}

func getType(s string) donut.RecordType {
	switch s {
	case "A":
		return donut.A
	case "AAAA":
		return donut.AAAA
	case "CNAME":
		return donut.CNAME
	case "MX":
		return donut.MX
	case "NS":
		return donut.NS
	case "PTR":
		return donut.PTR
	case "SOA":
		return donut.SOA
	case "SRV":
		return donut.SRV
	case "TXT":
		return donut.TXT
	default:
		return donut.A
	}
}
