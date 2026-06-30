package provider

import (
	"errors"
	"testing"
)

func TestFindProvider(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	if p == nil {
		t.Fatal("FindProvider(cloudflare) returned nil")
	}
	if p.Name != "CloudFlare" {
		t.Errorf("Name = %q, want %q", p.Name, "CloudFlare")
	}
}

func TestFindProvider_Unknown(t *testing.T) {
	p := FindProvider("nonexistent")
	if p != nil {
		t.Errorf("expected nil for unknown provider, got %v", p)
	}
}

func TestCloudFlare_Resolve_Default(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	primary, secondary, doh, err := p.Resolve("", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "1.1.1.1" {
		t.Errorf("primary = %q, want %q", primary, "1.1.1.1")
	}
	if secondary != "1.0.0.1" {
		t.Errorf("secondary = %q, want %q", secondary, "1.0.0.1")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty (useDOH=false)", doh)
	}
}

func TestCloudFlare_Resolve_StandardFlavor(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	primary, secondary, doh, err := p.Resolve("standard", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "1.1.1.1" {
		t.Errorf("primary = %q, want %q", primary, "1.1.1.1")
	}
	if secondary != "1.0.0.1" {
		t.Errorf("secondary = %q, want %q", secondary, "1.0.0.1")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestCloudFlare_Resolve_MalwareFlavor(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	primary, secondary, doh, err := p.Resolve("malware", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "1.1.1.2" {
		t.Errorf("primary = %q, want %q", primary, "1.1.1.2")
	}
	if secondary != "1.0.0.2" {
		t.Errorf("secondary = %q, want %q", secondary, "1.0.0.2")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestCloudFlare_Resolve_AdultFlavor(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	primary, secondary, doh, err := p.Resolve("adult", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "1.1.1.3" {
		t.Errorf("primary = %q, want %q", primary, "1.1.1.3")
	}
	if secondary != "1.0.0.3" {
		t.Errorf("secondary = %q, want %q", secondary, "1.0.0.3")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestCloudFlare_Resolve_WithDOH(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	_, _, doh, err := p.Resolve("standard", true)
	if err != nil {
		t.Fatalf("Resolve() with DoH error = %v", err)
	}
	if doh != "https://cloudflare-dns.com/dns-query" {
		t.Errorf("doh = %q, want %q", doh, "https://cloudflare-dns.com/dns-query")
	}
}

func TestCloudFlare_Resolve_InvalidFlavor(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	_, _, _, err := p.Resolve("invalid-flavor", false)
	if err == nil {
		t.Fatal("expected error for invalid flavor, got nil")
	}
}

func TestGoogle_Resolve_Default(t *testing.T) {
	p := FindProvider(ProviderGoogle)
	primary, secondary, doh, err := p.Resolve("", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "8.8.8.8" {
		t.Errorf("primary = %q, want %q", primary, "8.8.8.8")
	}
	if secondary != "8.8.4.4" {
		t.Errorf("secondary = %q, want %q", secondary, "8.8.4.4")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestGoogle_Resolve_WithDOH(t *testing.T) {
	p := FindProvider(ProviderGoogle)
	_, _, doh, err := p.Resolve("", true)
	if err != nil {
		t.Fatalf("Resolve() with DoH error = %v", err)
	}
	if doh != "https://dns.google/dns-query" {
		t.Errorf("doh = %q, want %q", doh, "https://dns.google/dns-query")
	}
}

func TestGoogle_Resolve_FlavorUnsupported(t *testing.T) {
	p := FindProvider(ProviderGoogle)
	// Google has no flavors, so any flavor should error
	_, _, _, err := p.Resolve("standard", false)
	if err == nil {
		t.Fatal("expected error for flavor on Google provider, got nil")
	}
}

func TestQuad9_Resolve_Default(t *testing.T) {
	p := FindProvider(ProviderQuad9)
	primary, secondary, doh, err := p.Resolve("", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "9.9.9.9" {
		t.Errorf("primary = %q, want %q", primary, "9.9.9.9")
	}
	if secondary != "149.112.112.112" {
		t.Errorf("secondary = %q, want %q", secondary, "149.112.112.112")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestQuad9_Resolve_WithDOH(t *testing.T) {
	p := FindProvider(ProviderQuad9)
	_, _, doh, err := p.Resolve("", true)
	if err != nil {
		t.Fatalf("Resolve() with DoH error = %v", err)
	}
	if doh != "https://dns.quad9.net/dns-query" {
		t.Errorf("doh = %q, want %q", doh, "https://dns.quad9.net/dns-query")
	}
}

func TestQuad9_Resolve_FlavorUnsupported(t *testing.T) {
	p := FindProvider(ProviderQuad9)
	_, _, _, err := p.Resolve("standard", false)
	if err == nil {
		t.Fatal("expected error for flavor on Quad9 provider, got nil")
	}
}

func TestISP_Resolve_FromDHCP(t *testing.T) {
	p := FindProvider(ProviderISP)
	if p == nil {
		t.Fatal("FindProvider(isp) returned nil")
	}
	if !p.HasDynamicDNS() {
		t.Error("ISP should have dynamic DNS")
	}

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return []string{"192.168.1.1", "10.0.0.1"}, nil
	}
	defer func() { detectDHCPDNS = orig }()

	primary, secondary, doh, err := p.Resolve("", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "192.168.1.1" {
		t.Errorf("primary = %q, want %q", primary, "192.168.1.1")
	}
	if secondary != "10.0.0.1" {
		t.Errorf("secondary = %q, want %q", secondary, "10.0.0.1")
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestISP_Resolve_SingleDNS(t *testing.T) {
	p := FindProvider(ProviderISP)

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return []string{"192.168.1.254"}, nil
	}
	defer func() { detectDHCPDNS = orig }()

	primary, secondary, doh, err := p.Resolve("", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "192.168.1.254" {
		t.Errorf("primary = %q, want %q", primary, "192.168.1.254")
	}
	if secondary != "" {
		t.Errorf("secondary = %q, want empty", secondary)
	}
	if doh != "" {
		t.Errorf("doh = %q, want empty", doh)
	}
}

func TestISP_Resolve_DHCPError(t *testing.T) {
	p := FindProvider(ProviderISP)

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return nil, errors.New("no DHCP lease found")
	}
	defer func() { detectDHCPDNS = orig }()

	_, _, _, err := p.Resolve("", false)
	if err == nil {
		t.Fatal("expected error when DHCP detection fails, got nil")
	}
}

func TestISP_Resolve_FlavorIgnored(t *testing.T) {
	p := FindProvider(ProviderISP)

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return []string{"10.0.0.1"}, nil
	}
	defer func() { detectDHCPDNS = orig }()

	// ISP should ignore flavor and always use DHCP
	primary, _, _, err := p.Resolve("standard", false)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if primary != "10.0.0.1" {
		t.Errorf("primary = %q, want %q", primary, "10.0.0.1")
	}
}

func TestDescribeDNS_CloudFlare(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	desc := p.DescribeDNS("")
	if desc != "1.1.1.1 / 1.0.0.1" {
		t.Errorf("DescribeDNS = %q, want %q", desc, "1.1.1.1 / 1.0.0.1")
	}
}

func TestDescribeDNS_CloudFlareMalware(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	desc := p.DescribeDNS("malware")
	if desc != "1.1.1.2 / 1.0.0.2" {
		t.Errorf("DescribeDNS = %q, want %q", desc, "1.1.1.2 / 1.0.0.2")
	}
}

func TestDescribeDNS_Google(t *testing.T) {
	p := FindProvider(ProviderGoogle)
	desc := p.DescribeDNS("")
	if desc != "8.8.8.8 / 8.8.4.4" {
		t.Errorf("DescribeDNS = %q, want %q", desc, "8.8.8.8 / 8.8.4.4")
	}
}

func TestDescribeDNS_ISP(t *testing.T) {
	p := FindProvider(ProviderISP)

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return []string{"192.168.1.1"}, nil
	}
	defer func() { detectDHCPDNS = orig }()

	desc := p.DescribeDNS("")
	want := "DHCP: 192.168.1.1"
	if desc != want {
		t.Errorf("DescribeDNS = %q, want %q", desc, want)
	}
}

func TestDescribeDNS_ISP_Error(t *testing.T) {
	p := FindProvider(ProviderISP)

	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return nil, errors.New("no DHCP")
	}
	defer func() { detectDHCPDNS = orig }()

	desc := p.DescribeDNS("")
	if desc != "Unknown (DHCP detection failed)" {
		t.Errorf("DescribeDNS = %q, want %q", desc, "Unknown (DHCP detection failed)")
	}
}

func TestISPDescribeDNS(t *testing.T) {
	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return []string{"10.0.0.1", "10.0.0.2"}, nil
	}
	defer func() { detectDHCPDNS = orig }()

	desc, err := ISPDescribeDNS()
	if err != nil {
		t.Fatalf("ISPDescribeDNS() error = %v", err)
	}
	want := "DHCP: 10.0.0.1, 10.0.0.2"
	if desc != want {
		t.Errorf("ISPDescribeDNS() = %q, want %q", desc, want)
	}
}

func TestISPDescribeDNS_Error(t *testing.T) {
	orig := detectDHCPDNS
	detectDHCPDNS = func() ([]string, error) {
		return nil, errors.New("fail")
	}
	defer func() { detectDHCPDNS = orig }()

	_, err := ISPDescribeDNS()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCloudFlare_SupportsFlavors(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	if !p.SupportsFlavors() {
		t.Error("CloudFlare should support flavors")
	}
}

func TestGoogle_SupportsFlavors(t *testing.T) {
	p := FindProvider(ProviderGoogle)
	if p.SupportsFlavors() {
		t.Error("Google should not support flavors")
	}
}

func TestISP_SupportsFlavors(t *testing.T) {
	p := FindProvider(ProviderISP)
	if p.SupportsFlavors() {
		t.Error("ISP should not support flavors")
	}
}

func TestFindFlavor(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	f := p.FindFlavor("standard")
	if f == nil {
		t.Fatal("FindFlavor(standard) returned nil")
	}
	if f.PrimaryDNS != "1.1.1.1" {
		t.Errorf("flavor PrimaryDNS = %q, want %q", f.PrimaryDNS, "1.1.1.1")
	}
}

func TestFindFlavor_NotFound(t *testing.T) {
	p := FindProvider(ProviderCloudFlare)
	f := p.FindFlavor("nonexistent")
	if f != nil {
		t.Errorf("expected nil for nonexistent flavor, got %v", f)
	}
}

func TestAllProviders(t *testing.T) {
	providers := AllProviders()
	if len(providers) != 4 {
		t.Errorf("AllProviders() returned %d providers, want 4", len(providers))
	}

	types := make(map[ProviderType]bool)
	for _, p := range providers {
		types[p.Type] = true
	}
	if !types[ProviderISP] {
		t.Error("missing ISP provider")
	}
	if !types[ProviderCloudFlare] {
		t.Error("missing CloudFlare provider")
	}
	if !types[ProviderGoogle] {
		t.Error("missing Google provider")
	}
	if !types[ProviderQuad9] {
		t.Error("missing Quad9 provider")
	}
}
