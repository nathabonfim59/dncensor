package provider

import (
	"strings"
)

func NewISP() *DNSProvider {
	return &DNSProvider{
		Type: ProviderISP,
		Name: "ISP (restore original from DHCP)",
	}
}

func ISPDescribeDNS() (string, error) {
	ips, err := detectDHCPDNS()
	if err != nil {
		return "", err
	}
	return "DHCP: " + strings.Join(ips, ", "), nil
}
