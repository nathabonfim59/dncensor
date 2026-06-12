package stack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type nmStack struct{}

func (s *nmStack) Type() StackType {
	return StackNetworkManager
}

func (s *nmStack) Detect() bool {
	path, err := exec.LookPath("nmcli")
	if err != nil {
		return false
	}
	out, err := exec.Command(path, "-t", "-f", "RUNNING", "general").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "running"
}

func (s *nmStack) activeConnection() (string, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "NAME,DEVICE,TYPE", "con", "show", "--active").Output()
	if err != nil {
		return "", fmt.Errorf("nmcli: %w", err)
	}

	// Prefer ethernet over wifi
	var ethernet, wifi string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			name := parts[0]
			typ := parts[len(parts)-1]
			switch typ {
			case "ethernet":
				ethernet = name
			case "wifi":
				if wifi == "" {
					wifi = name
				}
			}
		}
	}
	if ethernet != "" {
		return ethernet, nil
	}
	if wifi != "" {
		return wifi, nil
	}
	return "", fmt.Errorf("no active ethernet or wifi connection found")
}

func (s *nmStack) CurrentDNS() (string, error) {
	conn, err := s.activeConnection()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("nmcli", "-s", "-f", "ipv4.dns,ipv4.dns-search", "con", "show", conn).Output()
	if err != nil {
		return "", fmt.Errorf("nmcli con show: %w", err)
	}
	return string(out), nil
}

func (s *nmStack) SetDNS(primary, secondary string) error {
	conn, err := s.activeConnection()
	if err != nil {
		return err
	}

	dnsVal := primary
	if secondary != "" {
		dnsVal += " " + secondary
	}

	if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.dns", dnsVal).Run(); err != nil {
		return fmt.Errorf("nmcli con mod ipv4.dns: %w", err)
	}
	if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.ignore-auto-dns", "yes").Run(); err != nil {
		return fmt.Errorf("nmcli con mod ignore-auto-dns: %w", err)
	}
	if err := exec.Command("nmcli", "con", "up", conn).Run(); err != nil {
		return fmt.Errorf("nmcli con up: %w", err)
	}
	return nil
}

func (s *nmStack) SetDOH(endpoint string) error {
	return fmt.Errorf("NetworkManager does not support DoH natively")
}

func (s *nmStack) Backup(backupDir string) error {
	conn, err := s.activeConnection()
	if err != nil {
		return err
	}

	out, err := exec.Command("nmcli", "-s", "-f", "ipv4.dns,ipv4.dns-search", "con", "show", conn).Output()
	if err != nil {
		return fmt.Errorf("backup nmcli: %w", err)
	}

	path := filepath.Join(backupDir, fmt.Sprintf("original-%s.txt", s.Type()))
	return os.WriteFile(path, out, 0600)
}

func (s *nmStack) Restore(backupDir string) error {
	conn, err := s.activeConnection()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(backupDir, fmt.Sprintf("original-%s.txt", s.Type())))
	if err != nil {
		return fmt.Errorf("read backup: %w", err)
	}

	var dnsVals []string
	var searchDomains []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "ipv4.dns:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "ipv4.dns:"))
			if val != "" {
				dnsVals = strings.Fields(val)
			}
		}
		if strings.HasPrefix(line, "ipv4.dns-search:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "ipv4.dns-search:"))
			if val != "" {
				searchDomains = strings.Fields(val)
			}
		}
	}

	if len(dnsVals) == 0 {
		if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.dns", "").Run(); err != nil {
			return fmt.Errorf("nmcli con mod clear dns: %w", err)
		}
		if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.ignore-auto-dns", "no").Run(); err != nil {
			return fmt.Errorf("nmcli con mod ignore-auto-dns: %w", err)
		}
	} else {
		dnsStr := strings.Join(dnsVals, " ")
		if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.dns", dnsStr).Run(); err != nil {
			return fmt.Errorf("nmcli con mod restore dns: %w", err)
		}
		if err := exec.Command("nmcli", "con", "mod", conn, "ipv4.ignore-auto-dns", "yes").Run(); err != nil {
			return fmt.Errorf("nmcli con mod ignore-auto-dns: %w", err)
		}
	}

	if len(searchDomains) > 0 {
		searchStr := strings.Join(searchDomains, " ")
		exec.Command("nmcli", "con", "mod", conn, "ipv4.dns-search", searchStr).Run()
	}

	if err := exec.Command("nmcli", "con", "up", conn).Run(); err != nil {
		return fmt.Errorf("nmcli con up: %w", err)
	}

	return nil
}

func (s *nmStack) RequiresRoot() bool {
	return true
}
