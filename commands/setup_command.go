package commands

import (
	"bytes"
	"html/template"
	"os"
	"path"

	"embed"

	"github.com/spf13/cobra"
)

//go:embed template
var templates embed.FS

func setup() *cobra.Command {
	cmd := &cobra.Command{
		Use: "setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			corefileTemplateFile, err := templates.ReadFile("template/coredns/Corefile")
			if err != nil {
				return err
			}

			corefileTemplate := template.New("corefile")
			_, err = corefileTemplate.Parse(string(corefileTemplateFile))
			if err != nil {
				return err
			}

			ltld, _ := cmd.Flags().GetString("ltld")
			fallbackDNSes, _ := cmd.Flags().GetStringArray("fallback-dnses")

			buf := bytes.NewBuffer(nil)
			corefileTemplate.Execute(buf, map[string]any{
				"LTLD":          ltld,
				"FallbackDNSes": fallbackDNSes,
				"Name":          "{{ .Name }}",
			})

			traefikConfigFile, err := templates.ReadFile("template/traefik/traefik.yml")
			if err != nil {
				return err
			}

			traefikDynamicConfigFile, err := templates.ReadFile("template/traefik/dynamic.yml")
			if err != nil {
				return err
			}

			dockerComposeFile, err := templates.ReadFile("template/docker/compose.yml")
			if err != nil {
				return err
			}

			os.WriteFile(
				path.Join(ghostHome, "Corefile"),
				buf.Bytes(),
				0600,
			)

			os.WriteFile(
				path.Join(ghostHome, "config"),
				[]byte(ltld),
				0600,
			)

			os.WriteFile(
				path.Join(ghostHome, "traefik.yml"),
				traefikConfigFile,
				0600,
			)

			os.WriteFile(
				path.Join(ghostHome, "dynamic.yml"),
				traefikDynamicConfigFile,
				0600,
			)

			os.WriteFile(
				path.Join(ghostHome, "compose.yml"),
				dockerComposeFile,
				0600,
			)

			err = applyDNSConfig()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("ltld", "ghost", "Local Top-Level Domain to setup")
	cmd.Flags().StringArray("fallback-dnses", []string{"1.1.1.1", "8.8.8.8"}, "Fallback DNSes")
	return cmd
}
