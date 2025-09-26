//go:build darwin

package commands

import (
	"fmt"
	"os/exec"
)

func applyDNSConfig() error {
	script := `networksetup -listallnetworkservices | tail +2 | while IFS= read -r service; do
    	networksetup -setdnsservers "$service" 127.0.0.1 1.1.1.1		
	done`

	cmd := exec.Command("/bin/sh", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply DNS config: %v:\n%s", err, string(out))
	}
	return nil
}
