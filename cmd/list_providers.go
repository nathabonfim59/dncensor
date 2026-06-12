package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathabonfim59/dncensor/internal/provider"
)

var listJSONFlag bool

var listProvidersCmd = &cobra.Command{
	Use:   "list-providers",
	Short: "List available DNS providers and flavors",
	Long:  `List all available DNS providers and their flavors.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		providers := provider.AllProviders()

		if listJSONFlag {
			type flavorOut struct {
				Name        string `json:"name"`
				Display     string `json:"display"`
				PrimaryDNS  string `json:"primary_dns"`
				SecondaryDNS string `json:"secondary_dns"`
				DoHEndpoint string `json:"doh_endpoint,omitempty"`
			}
			type providerOut struct {
				Name         string       `json:"name"`
				Type         string       `json:"type"`
				Flavors      []flavorOut  `json:"flavors,omitempty"`
			}

			var out []providerOut
			for _, p := range providers {
				po := providerOut{
					Name: p.Name,
					Type: string(p.Type),
				}
				for _, f := range p.Flavors {
					po.Flavors = append(po.Flavors, flavorOut{
						Name:         f.FlavorName,
						Display:      f.Display,
						PrimaryDNS:   f.PrimaryDNS,
						SecondaryDNS: f.SecondaryDNS,
						DoHEndpoint:  f.DOHEndpoint,
					})
				}
				out = append(out, po)
			}

			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			for _, p := range providers {
				fmt.Printf("  %s (%s)\n", p.Name, p.Type)
				if p.SupportsFlavors() {
					for _, f := range p.Flavors {
						fmt.Printf("    - %s (%s)\n", f.Display, f.FlavorName)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	listProvidersCmd.Flags().BoolVarP(&listJSONFlag, "json", "j", false, "Output in JSON format")
}
