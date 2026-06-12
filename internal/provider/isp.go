package provider

func NewISP() *DNSProvider {
	return &DNSProvider{
		Type: ProviderISP,
		Name: "ISP (restore original)",
	}
}
