package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/nathabonfim59/dncensor/internal/stack"
	"github.com/nathabonfim59/dncensor/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "dncensor",
	Short: "Do Not Censor — DNS provider switcher for Linux",
	Long: `dncensor detects your DNS stack and lets you switch between
DNS providers (ISP, CloudFlare, Google) via an interactive TUI or CLI.

Examples:
  dncensor                           # Launch TUI
  dncensor set -p cloudflare         # Set CloudFlare DNS
  dncensor set -p cloudflare -f malware --doh
  dncensor current                   # Show current DNS
  dncensor list-providers            # List available providers
  dncensor backup create -n "before" # Save a named backup
  dncensor backup list               # List saved backups
  dncensor backup restore <name>     # Restore from backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireRoot("")

		if err := initConfig(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		s := stack.Detect()
		if s == nil {
			return fmt.Errorf("no supported DNS stack found")
		}

		m := tui.NewModel(s)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(listProvidersCmd)
	rootCmd.AddCommand(backupCmd)
}
