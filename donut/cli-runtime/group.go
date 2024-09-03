package cliruntime // import "github.com/tomasbasham/donut/cli-runtime"

import "github.com/spf13/cobra"

type CommandGroup struct {
	Title    string
	Commands []*cobra.Command
}

type CommandGroups []CommandGroup

func (g CommandGroups) Add(cmd *cobra.Command) {
	for _, group := range g {
		registerCommandGroup(cmd, group)
	}
}

func registerCommandGroup(cmd *cobra.Command, cg CommandGroup) {
	if len(cg.Title) == 0 {
		panic("CommandGroup requires a name")
	}

	group := &cobra.Group{
		ID:    cg.Title,
		Title: cg.Title,
	}

	for _, cc := range cg.Commands {
		cc.GroupID = cg.Title
		cmd.AddCommand(cc)
	}

	cmd.AddGroup(group)
}
