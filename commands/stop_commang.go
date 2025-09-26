package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func stopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stops ghost containers",
		Long:    "This command stops the ghost containers.",
		Example: "ghost start",
		RunE: func(cmd *cobra.Command, args []string) error {
			terminalCmd := exec.Command("docker",
				"compose", "-f", fmt.Sprintf("%s/compose.yml", ghostHome), "down")
			terminalCmd.Env = []string{
				fmt.Sprintf("GHOST_HOME=%s", ghostHome),
			}
			out, err := terminalCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to stop containers: %v:\n%s", err, string(out))
			}
			return nil
		},
	}
	return cmd
}
