package dns

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"
)

type ApplyConfig struct {
	Provider   *provider.DNSProvider
	FlavorName string
	UseDOH     bool
}

type ApplyResult struct {
	Success bool
	Message string
	ApplyConfig
}

func Apply(s stack.Stack, cfg ApplyConfig, backupDir string) ApplyResult {
	if cfg.Provider.Type == provider.ProviderISP {
		return restoreISP(s, backupDir)
	}

	primary, secondary, dohEndpoint, err := cfg.Provider.Resolve(cfg.FlavorName, cfg.UseDOH)
	if err != nil {
		return ApplyResult{
			Success:     false,
			Message:     fmt.Sprintf("Failed to resolve DNS config: %s", err),
			ApplyConfig: cfg,
		}
	}

	originalPath := filepath.Join(backupDir, fmt.Sprintf("original-%s.txt", s.Type()))
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		if err := s.Backup(backupDir); err != nil {
			return ApplyResult{
				Success:     false,
				Message:     fmt.Sprintf("Failed to backup original DNS config: %s", err),
				ApplyConfig: cfg,
			}
		}
	}

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
	originalPath := filepath.Join(backupDir, fmt.Sprintf("original-%s.txt", s.Type()))
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		return ApplyResult{
			Success: false,
			Message: "No original ISP backup found. Apply a non-ISP provider first to create a backup.",
		}
	}

	if err := s.Restore(backupDir); err != nil {
		return ApplyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to restore ISP settings: %s", err),
		}
	}

	return ApplyResult{
		Success: true,
		Message: "ISP settings restored from backup",
	}
}
