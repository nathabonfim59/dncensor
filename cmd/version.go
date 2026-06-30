package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionInfo = "dev"
	commitInfo  = "none"
	dateInfo    = "unknown"
)

func SetVersionInfo(version, commit, date string) {
	versionInfo = version
	commitInfo = commit
	dateInfo = date
	rootCmd.Version = version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dncensor %s\n", versionInfo)
		fmt.Printf("commit: %s\n", commitInfo)
		fmt.Printf("built:  %s\n", dateInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
