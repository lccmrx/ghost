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
			config, err := os.ReadFile(
				path.Join(ghostHome, "config"),
			)
			if err != nil || len(config) == 0 {
				return fmt.Errorf("error reading config file\n\n>> make sure you've run `ghost setup` before <<")
			}

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

			fmt.Println("Execute the following command on your terminal:")
			fmt.Printf("\n>> sudo sed '/%s.%s #ghost/d' /etc/hosts\n\n", args[0], config)

			return nil
		},
	}
	return cmd
}
