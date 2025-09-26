package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/lccmrx/ghost/commands/template/traefik"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func removeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [name]",
		Short:   "Remove a ghost service configuration",
		Long:    "This command removes a service from the ghost's configuration.",
		Example: "ghost remove example",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			traefikDynamicConfigFile, err := os.ReadFile(
				path.Join(ghostHome, "dynamic.yml"),
			)
			if err != nil {
				return err
			}

			var data traefik.DynamicConfig

			err = yaml.Unmarshal(traefikDynamicConfigFile, &data)
			if err != nil {
				return fmt.Errorf("error unmarshalling `dynamic.yml` file: %w", err)
			}

			delete(data.HTTP.Routers, fmt.Sprintf("%s-router", args[0]))
			delete(data.HTTP.Services, fmt.Sprintf("%s-service", args[0]))
			dataOut, _ := yaml.Marshal(data)

			os.WriteFile(
				path.Join(ghostHome, "dynamic.yml"),
				dataOut,
				0600,
			)

			return nil
		},
	}
	return cmd
}
