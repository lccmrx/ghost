package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/lccmrx/ghost/commands/template/traefik"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func addCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name] [port]",
		Short:   "Adds new service configuration",
		Long:    "This command adds a new service to the ghost's configuration.",
		Args:    cobra.ExactArgs(2),
		Example: "ghost add example 8000",
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

			if data.HTTP.Routers == nil {
				data.HTTP.Routers = make(map[string]traefik.RouterConfig)
			}

			data.HTTP.Routers[fmt.Sprintf("%s-router", args[0])] = traefik.RouterConfig{
				Rule:        fmt.Sprintf("Host(`%s.%s`)", args[0], config),
				Service:     fmt.Sprintf("%s-service", args[0]),
				EntryPoints: []string{"web"},
			}

			data.HTTP.Services[fmt.Sprintf("%s-service", args[0])] = traefik.ServiceConfig{
				LoadBalancer: traefik.LoadBalancerConfig{
					Servers: []traefik.ServerConfig{
						{
							URL: fmt.Sprintf("http://host.docker.internal:%s", args[1]),
						},
					},
				},
			}

			dataOut, _ := yaml.Marshal(data)

			os.WriteFile(
				path.Join(ghostHome, "dynamic.yml"),
				dataOut,
				0600,
			)

			fmt.Println("Execute the following command on your terminal:")
			fmt.Printf("\n>> echo -n '127.0.0.1 %s.%s #ghost\\n' | sudo tee -a /etc/hosts\n\n", args[0], config)

			return nil
		},
	}
	return cmd
}
