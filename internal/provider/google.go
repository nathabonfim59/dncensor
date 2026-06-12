package provider

func NewGoogle() *DNSProvider {
	return &DNSProvider{
		Type:         ProviderGoogle,
		Name:         "Google",
		PrimaryDNS:   "8.8.8.8",
		SecondaryDNS: "8.8.4.4",
		DOHEndpoint:  "https://dns.google/dns-query",
	}
}
