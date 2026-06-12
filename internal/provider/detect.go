package provider

import (
	"strings"

	"github.com/nathabonfim59/dncensor/internal/stack"
)

func DetectCurrentProvider(s stack.Stack) (provider *DNSProvider, flavor *DNSFlavor, err error) {
	raw, err := s.CurrentDNS()
	if err != nil {
		return nil, nil, err
	}

	var ips []string
	switch s.Type() {
	case stack.StackSystemdResolved:
		ips = parseResolvectlStatus(raw)
	case stack.StackNetworkManager:
		ips = parseNmcliDNS(raw)
	case stack.StackResolvConf:
		ips = parseResolvConf(raw)
	}

	if len(ips) == 0 {
		return nil, nil, nil
	}

	p, f := matchProvider(ips)
	return p, f, nil
}

func parseResolvectlStatus(output string) []string {
	seen := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "DNS Servers:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "DNS Servers:"))
		for _, part := range strings.Fields(rest) {
			part = strings.Trim(part, ",")
			if strings.Contains(part, ".") && !seen[part] {
				seen[part] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for ip := range seen {
		out = append(out, ip)
	}
	return out
}

func parseNmcliDNS(output string) []string {
	seen := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "ipv4.dns:") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "ipv4.dns:"))
		for _, part := range strings.Split(rest, ",") {
			part = strings.TrimSpace(part)
			if strings.Contains(part, ".") && !seen[part] {
				seen[part] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for ip := range seen {
		out = append(out, ip)
	}
	return out
}

func parseResolvConf(output string) []string {
	seen := map[string]bool{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "nameserver ") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ip := parts[1]
			if strings.Contains(ip, ".") && !seen[ip] {
				seen[ip] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for ip := range seen {
		out = append(out, ip)
	}
	return out
}

func matchProvider(ips []string) (provider *DNSProvider, flavor *DNSFlavor) {
	ipSet := make(map[string]bool, len(ips))
	for _, ip := range ips {
		ipSet[ip] = true
	}

	for _, p := range AllProviders() {
		if p.Type == ProviderISP {
			continue
		}

		for _, f := range p.Flavors {
			if ipSet[f.PrimaryDNS] && ipSet[f.SecondaryDNS] {
				return p, &DNSFlavor{
					Type:         f.Type,
					FlavorName:   f.FlavorName,
					Display:      f.Display,
					Description:  f.Description,
					PrimaryDNS:   f.PrimaryDNS,
					SecondaryDNS: f.SecondaryDNS,
					DOHEndpoint:  f.DOHEndpoint,
				}
			}
		}

		if ipSet[p.PrimaryDNS] && ipSet[p.SecondaryDNS] {
			return p, nil
		}
	}

	return nil, nil
}
