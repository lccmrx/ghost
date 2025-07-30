package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	dockerClient *client.Client
)

var (
	ghostConfigDir = os.Getenv("HOME") + "/.ghost"
)

func init() {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(ghostConfigDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/traefik.yaml", ghostConfigDir), []byte(`
entryPoints:
  web:
    address: ":80"

providers:
  file:
    filename: "/etc/traefik/dynamic.yaml"
    watch   : true
`), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func main() {
	rootCmd := cobra.Command{
		Use:   "ghost",
		Short: "ghost is a tool for managing Local Hosts on your host machine",
		Long:  "ghost is a command-line tool that allows you to manage your local hosts easily. The goal is to facilitate running multiple hosts and access them by local domain names.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:     "add",
			Short:   "Manage ghost configuration",
			Long:    "This command allows you to add a new host to the ghost configuration.",
			Args:    cobra.ExactArgs(2),
			Example: "ghost add example.com 8000",
			Run: func(cmd *cobra.Command, args []string) {
				file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Error opening /etc/hosts:", err)
					return
				}
				defer file.Close()

				_, err = file.WriteString(fmt.Sprintf("\n127.0.0.1\t%s #ghost\n", args[0]))
				if err != nil {
					fmt.Println("Error writing to /etc/hosts:", err)
					return
				}

				dynamicFile, err := os.ReadFile(fmt.Sprintf("%s/dynamic.yaml", ghostConfigDir))
				if err != nil {
					fmt.Println("Error opening dynamic.yaml:", err)
					return
				}

				var data DynamicConfig

				err = yaml.Unmarshal(dynamicFile, &data)
				if err != nil {
					fmt.Println("Error unmarshalling YAML:", err)
					return
				}

				data.HTTP.Routers[fmt.Sprintf("%s-router", args[0])] = RouterConfig{
					Rule:        fmt.Sprintf("Host(`%s`)", args[0]),
					Service:     fmt.Sprintf("%s-service", args[0]),
					EntryPoints: []string{"web"},
				}

				data.HTTP.Services[fmt.Sprintf("%s-service", args[0])] = ServiceConfig{
					LoadBalancer: LoadBalancerConfig{
						Servers: []ServerConfig{
							{
								URL: fmt.Sprintf("http://host.docker.internal:%s", args[1]),
							},
						},
					},
				}

				dataOut, _ := yaml.Marshal(data)

				os.WriteFile(fmt.Sprintf("%s/dynamic.yaml", ghostConfigDir), []byte(dataOut), 0644)
			},
		},
	)

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Starts the traefik service via Docker",
			Long:  "This command allows you to start the traefik service using Docker. It is necessary for managing local hosts.",
			Run: func(cmd *cobra.Command, args []string) {

				ctx := cmd.Context()
				runningContainers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
					Filters: filters.NewArgs(filters.Arg("label", "app=ghost")),
				})
				if err != nil {
					fmt.Println("Error fetching running containers:", err)
					return
				}

				if len(runningContainers) > 0 && runningContainers[0].State == "running" {
					fmt.Println("Traefik service is already running.")
					return
				}

				fmt.Println("Starting traefik service...")

				_, err = dockerClient.ImagePull(ctx, "traefik:latest", image.PullOptions{})
				if err != nil {
					fmt.Println("Error pulling traefik image:", err)
					return
				}

				traefikContainer, err := dockerClient.ContainerCreate(ctx, &container.Config{
					Image: "traefik:latest",
					ExposedPorts: nat.PortSet{
						"80/tcp": {},
					},
					Volumes: map[string]struct{}{
						"/etc/traefik": {},
					},
					Labels: map[string]string{
						"app":            "ghost",
						"traefik.enable": "true",
					},
					Cmd: []string{
						"--configFile=/etc/traefik/traefik.yml",
					},
				}, &container.HostConfig{
					PortBindings: nat.PortMap{
						"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "80"}},
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: ghostConfigDir,
							Target: "/etc/traefik",
						},
					},
				}, nil, nil, "ghost-traefik")
				if err != nil {
					fmt.Println("Error starting traefik service:", err)
					return
				}
				err = dockerClient.ContainerStart(ctx, traefikContainer.ID, container.StartOptions{})
				fmt.Println(err)
			},
		},
	)

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "clear",
			Short: "Clear all configs",
			Long:  "This command allows you to start the traefik service using Docker. It is necessary for managing local hosts.",
			Run: func(cmd *cobra.Command, args []string) {
				file, err := os.ReadFile("/etc/hosts")
				if err != nil {
					fmt.Println("Error opening /etc/hosts:", err)
					return
				}

				var hostsDataOut []string

				for _, line := range strings.Split(string(file), "\n") {
					line = strings.TrimSpace(line)
					if !strings.HasSuffix(line, " #ghost") {
						hostsDataOut = append(hostsDataOut, line)
					}
				}

				var data = DynamicConfig{}

				dataOut, _ := yaml.Marshal(data)

				os.WriteFile("/etc/hosts", []byte(strings.Join(hostsDataOut, "\n")), 0644)

				os.WriteFile(fmt.Sprintf("%s/dynamic.yaml", ghostConfigDir), []byte(dataOut), 0644)
			},
		},
	)

	rootCmd.Execute()
}
