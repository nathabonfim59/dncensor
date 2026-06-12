package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nathabonfim59/dncensor/internal/config"
	"github.com/nathabonfim59/dncensor/internal/dns"
	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"
)

var (
	providerFlag string
	flavorFlag   string
	dohFlag      bool
	yesFlag      bool
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set DNS provider",
	Long: `Set the DNS provider to use.

Examples:
  dncensor set -p cloudflare
  dncensor set -p cloudflare -f malware
  dncensor set -p cloudflare -f malware --doh
  dncensor set -p google
  dncensor set -p isp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireRoot("set")

		if err := config.Init(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		s := stack.Detect()
		if s == nil {
			return fmt.Errorf("no supported DNS stack found")
		}

		p := provider.FindProvider(provider.ProviderType(providerFlag))
		if p == nil {
			return fmt.Errorf("unknown provider: %s (use: isp, cloudflare, google)", providerFlag)
		}

		if flavorFlag != "" && !p.SupportsFlavors() {
			return fmt.Errorf("provider %s does not support flavors", p.Name)
		}

		if flavorFlag != "" {
			if p.FindFlavor(flavorFlag) == nil {
				return fmt.Errorf("unknown flavor: %s", flavorFlag)
			}
		}

		cfg := dns.ApplyConfig{
			Provider:   p,
			FlavorName: flavorFlag,
			UseDOH:     dohFlag,
		}

		if !yesFlag {
			fmt.Printf("Apply DNS configuration:\n")
			fmt.Printf("  Provider: %s\n", p.Name)
			if flavorFlag != "" {
				flavor := p.FindFlavor(flavorFlag)
				if flavor != nil {
					fmt.Printf("  Flavor:   %s\n", flavor.Display)
				}
			}
			fmt.Printf("  DoH:      %v\n", dohFlag)
			fmt.Print("Continue? [y/N]: ")

			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		result := dns.Apply(s, cfg, config.BackupPath())
		if result.Success {
			fmt.Println(result.Message)
		} else {
			fmt.Fprintln(os.Stderr, result.Message)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	setCmd.Flags().StringVarP(&providerFlag, "provider", "p", "", "DNS provider (isp, cloudflare, google)")
	setCmd.Flags().StringVarP(&flavorFlag, "flavor", "f", "", "Provider flavor (e.g., standard, malware, adult)")
	setCmd.Flags().BoolVar(&dohFlag, "doh", false, "Enable DNS-over-HTTPS")
	setCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Skip confirmation prompt")
	setCmd.MarkFlagRequired("provider")
}
