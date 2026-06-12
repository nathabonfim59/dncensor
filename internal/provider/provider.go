package provider

import "fmt"

type ProviderType string

const (
	ProviderISP        ProviderType = "isp"
	ProviderCloudFlare ProviderType = "cloudflare"
	ProviderGoogle     ProviderType = "google"
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

func AllProviders() []*DNSProvider {
	return []*DNSProvider{
		NewISP(),
		NewCloudFlare(),
		NewGoogle(),
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
