package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"
)

var jsonFlag bool

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current DNS configuration",
	Long:  `Display the current DNS configuration detected from the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := stack.Detect()
		if s == nil {
			return fmt.Errorf("no supported DNS stack found")
		}

		dnsStr, err := s.CurrentDNS()
		if err != nil {
			return fmt.Errorf("get current DNS: %w", err)
		}

		detectedProvider, detectedFlavor, _ := provider.DetectCurrentProvider(s)

		if jsonFlag {
			data := map[string]interface{}{
				"stack":  string(s.Type()),
				"config": dnsStr,
			}
			if detectedProvider != nil {
				data["provider"] = string(detectedProvider.Type)
				if detectedFlavor != nil {
					data["flavor"] = detectedFlavor.FlavorName
				}
			}
			out, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(out))
		} else {
			fmt.Printf("DNS Stack: %s\n", s.Type())
			if detectedProvider != nil {
				line := fmt.Sprintf("Detected Provider: %s", detectedProvider.Name)
				if detectedFlavor != nil {
					line += fmt.Sprintf(" > %s", detectedFlavor.Display)
				}
				fmt.Println(line)
			} else {
				fmt.Println("Detected Provider: Unknown")
			}
			fmt.Println("---")
			fmt.Print(dnsStr)
		}

		return nil
	},
}

func init() {
	currentCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
}
