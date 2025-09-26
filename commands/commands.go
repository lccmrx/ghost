package commands

import (
	"iter"
	"os"

	"github.com/spf13/cobra"
)

var ghostHome = os.Getenv("GHOST_HOME")

var commands = []func() *cobra.Command{
	setup,
	startCommand,
	stopCommand,
	addCommand,
	removeCommand,
}

func All() iter.Seq[*cobra.Command] {
	return func(yield func(*cobra.Command) bool) {
		for _, commandConstructor := range commands {
			command := commandConstructor()
			command.PreRun = func(cmd *cobra.Command, args []string) {
				ghostHome = os.Getenv("GHOST_HOME")
			}
			if !yield(command) {
				return
			}
		}
	}
}
