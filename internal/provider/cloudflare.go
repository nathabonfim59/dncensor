package provider

func NewCloudFlare() *DNSProvider {
	return &DNSProvider{
		Type:       ProviderCloudFlare,
		Name:       "CloudFlare",
		PrimaryDNS: "1.1.1.1",
		SecondaryDNS: "1.0.0.1",
		DOHEndpoint: "https://cloudflare-dns.com/dns-query",
		Flavors: []DNSFlavor{
			{
				Type:         ProviderCloudFlare,
				FlavorName:   "standard",
				Display:      "Standard",
				Description:  "1.1.1.1 - Fast and private",
				PrimaryDNS:   "1.1.1.1",
				SecondaryDNS: "1.0.0.1",
				DOHEndpoint:  "https://cloudflare-dns.com/dns-query",
			},
			{
				Type:         ProviderCloudFlare,
				FlavorName:   "malware",
				Display:      "Malware Blocking",
				Description:  "1.1.1.2 - Blocks malware",
				PrimaryDNS:   "1.1.1.2",
				SecondaryDNS: "1.0.0.2",
				DOHEndpoint:  "https://cloudflare-dns.com/dns-query",
			},
			{
				Type:         ProviderCloudFlare,
				FlavorName:   "adult",
				Display:      "Malware + Adult Content",
				Description:  "1.1.1.3 - Blocks malware and adult content",
				PrimaryDNS:   "1.1.1.3",
				SecondaryDNS: "1.0.0.3",
				DOHEndpoint:  "https://cloudflare-dns.com/dns-query",
			},
		},
	}
}
