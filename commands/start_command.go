package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func startCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Starts ghost containers",
		Long:    "This command starts the local DNS server and the Traefik proxy.",
		Example: "ghost start",
		RunE: func(cmd *cobra.Command, args []string) error {
			terminalCmd := exec.Command("docker",
				"compose", "-f", fmt.Sprintf("%s/compose.yml", ghostHome), "up", "-d")
			terminalCmd.Env = []string{
				fmt.Sprintf("GHOST_HOME=%s", ghostHome),
			}
			out, err := terminalCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to start containers: %v:\n%s", err, string(out))
			}
			fmt.Println(string(out))
			return nil
		},
	}
	return cmd
}
