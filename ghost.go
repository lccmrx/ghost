package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/lccmrx/ghost/commands"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "ghost",
		Short: "ghost is a tool for managing Local Hosts on your host machine",
		Long: `ghost is a command-line tool that allows you to manage your local hosts easily.
The goal is to facilitate running multiple hosts and access them by local domain names.
So if you are developing a multi services running multiple ports, this CLI is made for you!`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			os.Setenv("GHOST_HOME", path.Join(os.Getenv("HOME"), ".ghost"))
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			os.Unsetenv("GHOST_HOME")
		},
	}

	for command := range commands.All() {
		rootCmd.AddCommand(command)
	}

	rootCmd.Execute()
}
