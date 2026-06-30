package provider

func NewQuad9() *DNSProvider {
	return &DNSProvider{
		Type:         ProviderQuad9,
		Name:         "Quad9",
		PrimaryDNS:   "9.9.9.9",
		SecondaryDNS: "149.112.112.112",
		DOHEndpoint:  "https://dns.quad9.net/dns-query",
	}
}
