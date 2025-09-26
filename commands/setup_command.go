package commands

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
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
			fmt.Printf("using `%s` as GHOST_HOME\n", ghostHome)

			err := os.Mkdir(ghostHome, 0755)
			if err != nil && !errors.Is(err, fs.ErrExist) {
				return fmt.Errorf("failed to create `GHOST_HOME`: %w", err)
			}

			corefileTemplateFile, err := templates.ReadFile("template/coredns/Corefile")
			if err != nil {
				return fmt.Errorf("failed to read template file: %w", err)
			}

			corefileTemplate := template.New("corefile")
			_, err = corefileTemplate.Parse(string(corefileTemplateFile))
			if err != nil {
				return fmt.Errorf("failed to parse template file: %w", err)
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
				return fmt.Errorf("failed to read template file: %w", err)
			}

			traefikDynamicConfigFile, err := templates.ReadFile("template/traefik/dynamic.yml")
			if err != nil {
				return fmt.Errorf("failed to read template file: %w", err)
			}

			dockerComposeFile, err := templates.ReadFile("template/docker/compose.yml")
			if err != nil {
				return fmt.Errorf("failed to read template file: %w", err)
			}

			err = os.WriteFile(
				path.Join(ghostHome, "Corefile"),
				buf.Bytes(),
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			err = os.WriteFile(
				path.Join(ghostHome, "config"),
				[]byte(ltld),
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			err = os.WriteFile(
				path.Join(ghostHome, "traefik.yml"),
				traefikConfigFile,
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			err = os.WriteFile(
				path.Join(ghostHome, "dynamic.yml"),
				traefikDynamicConfigFile,
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			err = os.WriteFile(
				path.Join(ghostHome, "compose.yml"),
				dockerComposeFile,
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

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
