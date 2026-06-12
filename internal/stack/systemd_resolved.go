package stack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

func (s *systemdResolvedStack) Backup(backupDir string) error {
	ts := time.Now().Format("2006-01-02T15-04-05")

	statusOut, err := exec.Command("resolvectl", "status").Output()
	if err != nil {
		return fmt.Errorf("backup resolvectl status: %w", err)
	}
	if err := os.WriteFile(filepath.Join(backupDir, fmt.Sprintf("resolved-%s.txt", ts)), statusOut, 0600); err != nil {
		return err
	}

	dnsOut, err := exec.Command("resolvectl", "dns").Output()
	if err == nil {
		os.WriteFile(filepath.Join(backupDir, fmt.Sprintf("resolved-dns-%s.txt", ts)), dnsOut, 0600)
	}

	return nil
}

func (s *systemdResolvedStack) Restore(backupDir string) error {
	// For restore, we try to re-apply the previous DNS from backup files
	// This is best-effort; the full state is saved for manual recovery
	return fmt.Errorf("automatic restore for systemd-resolved not implemented; backup files are at %s", backupDir)
}

func (s *systemdResolvedStack) RequiresRoot() bool {
	return true
}
