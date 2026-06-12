package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nathabonfim59/dncensor/internal/backup"
	"github.com/nathabonfim59/dncensor/internal/config"
	"github.com/nathabonfim59/dncensor/internal/stack"
)

var backupName string
var backupYes bool

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage DNS configuration backups",
	Long: `Create, list, restore, and delete named backups of your DNS configuration.

Each backup is identified by a SHA256 hash and a user-provided name.
Restore by hash prefix or exact name.`,
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup of the current DNS config",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := stack.Detect()
		if s == nil {
			return fmt.Errorf("no supported DNS stack found")
		}

		if err := initConfig(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		name := backupName
		if name == "" {
			fmt.Print("Enter a name for this backup: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
			if name == "" {
				return fmt.Errorf("name cannot be empty")
			}
		}

		b, err := backup.Create(config.SnapshotsPath(), name, s)
		if err != nil {
			return fmt.Errorf("create backup: %w", err)
		}

		fmt.Printf("Backup created: %s (%s)\n", b.Hash[:12], b.Name)
		return nil
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		backups, err := backup.List(config.SnapshotsPath())
		if err != nil {
			return fmt.Errorf("list backups: %w", err)
		}

		if len(backups) == 0 {
			fmt.Println("No backups found.")
			return nil
		}

		fmt.Printf("%-12s  %-30s  %-20s  %s\n", "HASH", "NAME", "STACK", "CREATED")
		fmt.Println(strings.Repeat("-", 80))
		for _, b := range backups {
			fmt.Printf("%-12s  %-30s  %-20s  %s\n",
				b.Hash[:12], b.Name, b.StackType, b.CreatedAt.Format("2006-01-02 15:04"))
		}

		return nil
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <hash-or-name>",
	Short: "Restore DNS from a backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireRoot("backup restore")

		if err := initConfig(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		s := stack.Detect()
		if s == nil {
			return fmt.Errorf("no supported DNS stack found")
		}

		b, err := backup.Find(config.SnapshotsPath(), args[0])
		if err != nil {
			return fmt.Errorf("find backup: %w", err)
		}

		if !backupYes {
			fmt.Printf("Restore backup %s (%s) from %s?\n", b.Hash[:12], b.Name, b.StackType)
			fmt.Print("Continue? [y/N]: ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		if err := b.Restore(s); err != nil {
			return fmt.Errorf("restore backup: %w", err)
		}

		fmt.Printf("Restored backup: %s (%s)\n", b.Hash[:12], b.Name)
		return nil
	},
}

var backupDeleteCmd = &cobra.Command{
	Use:   "delete <hash-or-name>",
	Short: "Delete a DNS backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return fmt.Errorf("config init: %w", err)
		}

		b, err := backup.Find(config.SnapshotsPath(), args[0])
		if err != nil {
			return fmt.Errorf("find backup: %w", err)
		}

		if !backupYes {
			fmt.Printf("Delete backup %s (%s) from %s?\n", b.Hash[:12], b.Name, b.StackType)
			fmt.Print("Continue? [y/N]: ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		if err := backup.Delete(config.SnapshotsPath(), args[0]); err != nil {
			return fmt.Errorf("delete backup: %w", err)
		}

		fmt.Printf("Deleted backup: %s (%s)\n", b.Hash[:12], b.Name)
		return nil
	},
}

func init() {
	backupCreateCmd.Flags().StringVarP(&backupName, "name", "n", "", "Backup name (prompted if empty)")
	backupRestoreCmd.Flags().BoolVarP(&backupYes, "yes", "y", false, "Skip confirmation")
	backupDeleteCmd.Flags().BoolVarP(&backupYes, "yes", "y", false, "Skip confirmation")

	backupCmd.AddCommand(backupCreateCmd)
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupDeleteCmd)
}
