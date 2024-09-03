package main

import (
	cliruntime "github.com/tomasbasham/donut/cli-runtime"
	"github.com/tomasbasham/donut/internal/cmd"
)

func main() {
	command := cmd.NewRootCommand()
	if err := cliruntime.RunNoErrOutput(command); err != nil {
		panic(err)
	}
}
