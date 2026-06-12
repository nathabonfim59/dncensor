package stack

import (
	"fmt"
	"os/exec"
	"strings"
)

type systemdResolvedStack struct{}

func (s *systemdResolvedStack) Type() StackType {
	return StackSystemdResolved
}

func (s *systemdResolvedStack) Detect() bool {
	path, err := exec.LookPath("resolvectl")
	if err != nil {
		return false
	}
	cmd := exec.Command(path, "status")
	return cmd.Run() == nil
}

func (s *systemdResolvedStack) CurrentDNS() (string, error) {
	out, err := exec.Command("resolvectl", "status").Output()
	if err != nil {
		return "", fmt.Errorf("resolvectl status: %w", err)
	}
	return string(out), nil
}

func (s *systemdResolvedStack) interfaces() ([]string, error) {
	out, err := exec.Command("resolvectl", "status").Output()
	if err != nil {
		return nil, fmt.Errorf("resolvectl status: %w", err)
	}
	var ifaces []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Link ") {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) >= 3 {
				name := strings.Trim(parts[2], "()")
				if name != "" {
					ifaces = append(ifaces, name)
				}
			}
		}
	}
	return ifaces, nil
}

func (s *systemdResolvedStack) SetDNS(primary, secondary string) error {
	ifaces, err := s.interfaces()
	if err != nil {
		// Fallback: try lo
		ifaces = []string{"lo"}
	}

	for _, iface := range ifaces {
		args := []string{"dns", iface, primary}
		if secondary != "" {
			args = append(args, secondary)
		}
		if err := exec.Command("resolvectl", args...).Run(); err != nil {
			return fmt.Errorf("resolvectl dns %s: %w", iface, err)
		}
		if err := exec.Command("resolvectl", "domain", iface, "~.").Run(); err != nil {
			return fmt.Errorf("resolvectl domain %s: %w", iface, err)
		}
	}
	return nil
}

func (s *systemdResolvedStack) SetDOH(endpoint string) error {
	ifaces, err := s.interfaces()
	if err != nil {
		ifaces = []string{"lo"}
	}

	for _, iface := range ifaces {
		if err := exec.Command("resolvectl", "dns", iface, endpoint).Run(); err != nil {
			return fmt.Errorf("resolvectl dns %s (DoH): %w", iface, err)
		}
	}
	return nil
}

func (s *systemdResolvedStack) CaptureDNS() ([]byte, error) {
	out, err := exec.Command("resolvectl", "dns").Output()
	if err != nil {
		return nil, fmt.Errorf("capture resolvectl dns: %w", err)
	}
	return out, nil
}

func (s *systemdResolvedStack) ApplyDNS(content []byte) error {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Global:") {
			continue
		}
		if strings.HasPrefix(line, "Link ") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 2 {
				continue
			}
			ifacePart := strings.TrimPrefix(parts[0], "Link ")
			ifaceName := ""
			if idx := strings.Index(ifacePart, "("); idx != -1 {
				ifaceName = strings.TrimSuffix(ifacePart[idx+1:], ")")
			}
			if ifaceName == "" {
				continue
			}

			dnsIps := strings.Fields(strings.TrimSpace(parts[1]))
			if len(dnsIps) == 0 {
				continue
			}

			args := append([]string{"dns", ifaceName}, dnsIps...)
			if err := exec.Command("resolvectl", args...).Run(); err != nil {
				return fmt.Errorf("apply resolvectl dns %s: %w", ifaceName, err)
			}
			if err := exec.Command("resolvectl", "domain", ifaceName, "~.").Run(); err != nil {
				return fmt.Errorf("apply resolvectl domain %s: %w", ifaceName, err)
			}
		}
	}

	return nil
}

func (s *systemdResolvedStack) RequiresRoot() bool {
	return true
}
