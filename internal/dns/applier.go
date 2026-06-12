package dns

import (
	"fmt"

	"github.com/nathabonfim59/dncensor/internal/backup"
	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"
)

type ApplyConfig struct {
	Provider    *provider.DNSProvider
	FlavorName  string
	UseDOH      bool
}

type ApplyResult struct {
	Success bool
	Message string
	ApplyConfig
}

func Apply(s stack.Stack, cfg ApplyConfig, backupDir string) ApplyResult {
	// ISP = restore from backup
	if cfg.Provider.Type == provider.ProviderISP {
		return restoreISP(s, backupDir)
	}

	// Resolve DNS servers
	primary, secondary, dohEndpoint, err := cfg.Provider.Resolve(cfg.FlavorName, cfg.UseDOH)
	if err != nil {
		return ApplyResult{
			Success:     false,
			Message:     fmt.Sprintf("Failed to resolve DNS config: %s", err),
			ApplyConfig: cfg,
		}
	}

	// Backup current config
	if err := s.Backup(backupDir); err != nil {
		return ApplyResult{
			Success:     false,
			Message:     fmt.Sprintf("Failed to backup current DNS config: %s", err),
			ApplyConfig: cfg,
		}
	}

	// Apply DNS servers
	if err := s.SetDNS(primary, secondary); err != nil {
		return ApplyResult{
			Success:     false,
			Message:     fmt.Sprintf("Failed to set DNS: %s", err),
			ApplyConfig: cfg,
		}
	}

	msg := fmt.Sprintf("DNS set to %s", cfg.Provider.Name)
	if cfg.FlavorName != "" {
		flavor := cfg.Provider.FindFlavor(cfg.FlavorName)
		if flavor != nil {
			msg += fmt.Sprintf(" (%s)", flavor.Display)
		}
	}

	// Apply DoH if requested
	if cfg.UseDOH && dohEndpoint != "" {
		if err := s.SetDOH(dohEndpoint); err != nil {
			msg += fmt.Sprintf(". DNS servers applied but DoH failed: %s", err)
		} else {
			msg += " with DoH enabled"
		}
	}

	return ApplyResult{
		Success:     true,
		Message:     msg,
		ApplyConfig: cfg,
	}
}

func restoreISP(s stack.Stack, backupDir string) ApplyResult {
	bm := backup.New(backupDir)

	record, err := bm.Latest(string(s.Type()))
	if err != nil {
		return ApplyResult{
			Success: false,
			Message: "No ISP backup found. Apply a non-ISP provider first to create a backup.",
		}
	}

	if err := s.Restore(record.BackupPath); err != nil {
		// Fallback: try to restore from specific stack restore method
		if err2 := s.Restore(backupDir); err2 != nil {
			return ApplyResult{
				Success: false,
				Message: fmt.Sprintf("Failed to restore ISP settings: %s", err2),
			}
		}
	}

	return ApplyResult{
		Success: true,
		Message: "ISP settings restored from backup",
	}
}
