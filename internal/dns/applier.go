package dns

import (
	"fmt"

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

func Apply(s stack.Stack, cfg ApplyConfig) ApplyResult {
	primary, secondary, dohEndpoint, err := cfg.Provider.Resolve(cfg.FlavorName, cfg.UseDOH)
	if err != nil {
		return ApplyResult{
			Success:     false,
			Message:     fmt.Sprintf("Failed to resolve DNS config: %s", err),
			ApplyConfig: cfg,
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
	} else if cfg.Provider.HasDynamicDNS() {
		msg += fmt.Sprintf(" (%s)", cfg.Provider.DescribeDNS(cfg.FlavorName))
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
