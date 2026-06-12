package dhcp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func DetectOriginalDNS() ([]string, error) {
	iface, err := defaultInterface()
	if err != nil {
		return nil, fmt.Errorf("find default interface: %w", err)
	}

	candidates := []func(string) ([]string, bool){
		parseDHCPCD,
		parseDHClient,
		func(_ string) ([]string, bool) { return parseSystemdNetworkd() },
	}

	for _, fn := range candidates {
		ips, ok := fn(iface)
		if ok && len(ips) > 0 {
			return ips, nil
		}
	}

	return nil, fmt.Errorf("no DHCP lease with DNS found for interface %s", iface)
}

func defaultInterface() (string, error) {
	out, err := exec.Command("ip", "route", "show", "default").Output()
	if err != nil {
		return "", fmt.Errorf("ip route: %w", err)
	}
	// "default via 192.168.0.1 dev wlan0 proto dhcp src 192.168.0.100 metric 600"
	parts := strings.Fields(string(out))
	for i, p := range parts {
		if p == "dev" && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("no default route found")
}

func parseDHCPCD(iface string) ([]string, bool) {
	paths := []string{
		filepath.Join("/var/lib/dhcpcd", fmt.Sprintf("lease-%s", iface)),
		filepath.Join("/var/lib/dhcpcd", fmt.Sprintf("dhcpcd-%s.lease", iface)),
		filepath.Join("/var/db/dhcpcd", fmt.Sprintf("lease-%s", iface)),
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if ips, ok := parseDHCPCDContent(data); ok {
			return ips, true
		}
	}

	return nil, false
}

func parseDHCPCDContent(data []byte) ([]string, bool) {
	re := regexp.MustCompile(`option domain_name_servers\s+(.+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) >= 2 {
		ips := strings.Fields(matches[1])
		if len(ips) > 0 {
			return ips, true
		}
	}
	return nil, false
}

func parseDHClient(iface string) ([]string, bool) {
	path := filepath.Join("/var/lib/dhcp", fmt.Sprintf("dhclient-%s.leases", iface))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return parseDHClientContent(data)
}

func parseDHClientContent(data []byte) ([]string, bool) {
	re := regexp.MustCompile(`option domain-name-servers\s+([^;]+);`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) >= 2 {
		raw := strings.TrimSpace(matches[1])
		ips := strings.FieldsFunc(raw, func(r rune) bool {
			return r == ',' || r == ' '
		})
		var cleaned []string
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				cleaned = append(cleaned, ip)
			}
		}
		if len(cleaned) > 0 {
			return cleaned, true
		}
	}

	return nil, false
}

func parseSystemdNetworkd() ([]string, bool) {
	entries, err := os.ReadDir("/run/systemd/netif/leases")
	if err != nil {
		return nil, false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/run/systemd/netif/leases", entry.Name()))
		if err != nil {
			continue
		}

		if ips, ok := parseSystemdNetworkdContent(data); ok {
			return ips, true
		}
	}

	return nil, false
}

func parseSystemdNetworkdContent(data []byte) ([]string, bool) {
	content := string(data)

	re := regexp.MustCompile(`(?m)^DNS=(.*)$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 2 {
		ips := strings.Fields(matches[1])
		if len(ips) > 0 {
			return ips, true
		}
	}

	return nil, false
}
