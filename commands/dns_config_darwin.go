//go:build darwin

package commands

import (
	"fmt"
	"os"
	"path"
)

func applyDNSConfig() error {
	config, err := os.ReadFile(
		path.Join(ghostHome, "config"),
	)
	if err != nil || len(config) == 0 {
		return fmt.Errorf("error reading config file\n\n>> make sure you've run `ghost setup` before <<")
	}
	fmt.Println("\nthese next steps will setup the resolver for your LTLD.\n>> these require `sudo` <<\nPlease run in your terminal!")

	fmt.Println("\n>>  sudo mkdir -p /etc/resolver")
	fmt.Printf(">>  echo \"nameserver 127.0.0.1\" | sudo tee /etc/resolver/%s > /dev/null\n", config)
	fmt.Println(">>  sudo dscacheutil -flushcache")
	fmt.Println(">>  sudo killall -HUP mDNSResponder\n")

	return nil
}
