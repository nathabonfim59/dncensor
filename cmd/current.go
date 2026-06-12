package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

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

		if jsonFlag {
			data := map[string]string{
				"stack":  string(s.Type()),
				"config": dnsStr,
			}
			out, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(out))
		} else {
			fmt.Printf("DNS Stack: %s\n", s.Type())
			fmt.Println("---")
			fmt.Print(dnsStr)
		}

		return nil
	},
}

func init() {
	currentCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output in JSON format")
}
