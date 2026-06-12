package dhcp

import (
	"reflect"
	"testing"
)

func TestParseDHCPCDContent(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []string
		wantOK  bool
	}{
		{
			name:   "single DNS",
			data:   `option domain_name_servers 192.168.1.1`,
			want:   []string{"192.168.1.1"},
			wantOK: true,
		},
		{
			name:   "multiple DNS",
			data:   `option domain_name_servers 192.168.1.1 8.8.8.8`,
			want:   []string{"192.168.1.1", "8.8.8.8"},
			wantOK: true,
		},
		{
			name:   "no DNS option",
			data:   `option routers 192.168.1.1`,
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty data",
			data:   "",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "DNS with trailing whitespace",
			data:   `option domain_name_servers 10.0.0.1 `,
			want:   []string{"10.0.0.1"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseDHCPCDContent([]byte(tt.data))
			if ok != tt.wantOK {
				t.Errorf("parseDHCPCDContent() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDHCPCDContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDHClientContent(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   []string
		wantOK bool
	}{
		{
			name:   "single DNS",
			data:   `option domain-name-servers 192.168.1.1;`,
			want:   []string{"192.168.1.1"},
			wantOK: true,
		},
		{
			name:   "multiple DNS comma-separated",
			data:   `option domain-name-servers 192.168.1.1, 8.8.8.8;`,
			want:   []string{"192.168.1.1", "8.8.8.8"},
			wantOK: true,
		},
		{
			name:   "multiple DNS space-separated",
			data:   `option domain-name-servers 192.168.1.1 8.8.8.8;`,
			want:   []string{"192.168.1.1", "8.8.8.8"},
			wantOK: true,
		},
		{
			name:   "no DNS option",
			data:   `option routers 192.168.1.1;`,
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty data",
			data:   "",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "DNS with extra whitespace",
			data:   `option domain-name-servers  10.0.0.1  ;`,
			want:   []string{"10.0.0.1"},
			wantOK: true,
		},
		{
			name: "real dhclient lease excerpt",
			data: `lease {
  interface "eth0";
  fixed-address 192.168.1.100;
  option domain-name-servers 1.1.1.1, 1.0.0.1;
  option routers 192.168.1.1;
}`,
			want:   []string{"1.1.1.1", "1.0.0.1"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseDHClientContent([]byte(tt.data))
			if ok != tt.wantOK {
				t.Errorf("parseDHClientContent() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDHClientContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSystemdNetworkdContent(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		want   []string
		wantOK bool
	}{
		{
			name:   "real systemd-networkd lease format",
			data:   "ADDRESS=192.168.0.100\nDNS=187.45.96.96 187.45.97.97\n",
			want:   []string{"187.45.96.96", "187.45.97.97"},
			wantOK: true,
		},
		{
			name:   "no DNS field",
			data:   "ADDRESS=192.168.0.100\nROUTER=192.168.0.1\n",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty DNS",
			data:   "ADDRESS=192.168.0.100\nDNS=\n",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty data",
			data:   "",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "legacy format with INTERFACE",
			data:   "INTERFACE=eth0\nDNS=1.1.1.1 1.0.0.1\n",
			want:   []string{"1.1.1.1", "1.0.0.1"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseSystemdNetworkdContent([]byte(tt.data))
			if ok != tt.wantOK {
				t.Errorf("parseSystemdNetworkdContent() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSystemdNetworkdContent() = %v, want %v", got, tt.want)
			}
		})
	}
}


