//go:build !dev

package cmd

import (
	"fmt"
	"os"
)

// requireRoot exits the process if not running as root.
// In dev builds (tag "dev"), this is a no-op.
func requireRoot(context string) {
	if os.Geteuid() != 0 {
		fmt.Printf("dncensor needs root to change DNS. Please re-run with: sudo dncensor %s\n", context)
		os.Exit(1)
	}
}
