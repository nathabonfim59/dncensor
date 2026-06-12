package provider

import (
	"strings"

	"github.com/nathabonfim59/dncensor/internal/dhcp"
)

func NewISP() *DNSProvider {
	return &DNSProvider{
		Type: ProviderISP,
		Name: "ISP (restore original from DHCP)",
	}
}

func ISPDescribeDNS() (string, error) {
	ips, err := dhcp.DetectOriginalDNS()
	if err != nil {
		return "", err
	}
	return "DHCP: " + strings.Join(ips, ", "), nil
}
