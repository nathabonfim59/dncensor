package provider

import (
	"fmt"
	"strings"

	"github.com/nathabonfim59/dncensor/internal/dhcp"
)

// overridable in tests
var detectDHCPDNS = dhcp.DetectOriginalDNS

type ProviderType string

const (
	ProviderISP        ProviderType = "isp"
	ProviderCloudFlare ProviderType = "cloudflare"
	ProviderGoogle     ProviderType = "google"
	ProviderQuad9      ProviderType = "quad9"
)

type DNSFlavor struct {
	Type         ProviderType
	FlavorName   string
	Display      string
	Description  string
	PrimaryDNS   string
	SecondaryDNS string
	DOHEndpoint  string
}

type DNSProvider struct {
	Type         ProviderType
	Name         string
	PrimaryDNS   string
	SecondaryDNS string
	Flavors      []DNSFlavor
	DOHEndpoint  string
}

func (p *DNSProvider) SupportsFlavors() bool {
	return len(p.Flavors) > 0
}

func (p *DNSProvider) FindFlavor(name string) *DNSFlavor {
	for _, f := range p.Flavors {
		if f.FlavorName == name {
			return &f
		}
	}
	return nil
}

func (p *DNSProvider) Resolve(flavorName string, useDOH bool) (primary, secondary, dohEndpoint string, err error) {
	if p.Type == ProviderISP {
		ips, err := detectDHCPDNS()
		if err != nil {
			return "", "", "", fmt.Errorf("detect ISP DNS from DHCP: %w", err)
		}
		dohEndpoint = p.DOHEndpoint
		primary = ips[0]
		if len(ips) > 1 {
			secondary = ips[1]
		}
		if useDOH && dohEndpoint == "" {
			return primary, secondary, "", fmt.Errorf("DoH not supported for this provider/flavor")
		}
		if !useDOH {
			dohEndpoint = ""
		}
		return primary, secondary, dohEndpoint, nil
	}

	if flavorName != "" {
		flavor := p.FindFlavor(flavorName)
		if flavor == nil {
			return "", "", "", fmt.Errorf("flavor %q not found for provider %s", flavorName, p.Type)
		}
		primary = flavor.PrimaryDNS
		secondary = flavor.SecondaryDNS
		dohEndpoint = flavor.DOHEndpoint
	} else {
		primary = p.PrimaryDNS
		secondary = p.SecondaryDNS
		dohEndpoint = p.DOHEndpoint
	}

	if useDOH {
		if dohEndpoint == "" {
			return primary, secondary, "", fmt.Errorf("DoH not supported for this provider/flavor")
		}
	} else {
		dohEndpoint = ""
	}

	return primary, secondary, dohEndpoint, nil
}

func (p *DNSProvider) HasDynamicDNS() bool {
	return p.Type == ProviderISP
}

func (p *DNSProvider) DescribeDNS(flavorName string) string {
	if p.Type == ProviderISP {
		ips, err := detectDHCPDNS()
		if err != nil {
			return "Unknown (DHCP detection failed)"
		}
		return fmt.Sprintf("DHCP: %s", strings.Join(ips, ", "))
	}
	if flavorName != "" {
		flavor := p.FindFlavor(flavorName)
		if flavor != nil {
			return fmt.Sprintf("%s / %s", flavor.PrimaryDNS, flavor.SecondaryDNS)
		}
	}
	return fmt.Sprintf("%s / %s", p.PrimaryDNS, p.SecondaryDNS)
}

func AllProviders() []*DNSProvider {
	return []*DNSProvider{
		NewISP(),
		NewCloudFlare(),
		NewGoogle(),
		NewQuad9(),
	}
}

func FindProvider(typ ProviderType) *DNSProvider {
	for _, p := range AllProviders() {
		if p.Type == typ {
			return p
		}
	}
	return nil
}
