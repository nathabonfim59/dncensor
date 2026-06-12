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
		iface  string
		data   string
		want   []string
		wantOK bool
	}{
		{
			name:   "matching interface with DNS",
			iface:  "eth0",
			data:   "INTERFACE=eth0\nDNS=192.168.1.1 8.8.8.8\n",
			want:   []string{"192.168.1.1", "8.8.8.8"},
			wantOK: true,
		},
		{
			name:   "non-matching interface",
			iface:  "wlan0",
			data:   "INTERFACE=eth0\nDNS=192.168.1.1\n",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "no DNS field",
			iface:  "eth0",
			data:   "INTERFACE=eth0\nROUTES=...\n",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty DNS",
			iface:  "eth0",
			data:   "INTERFACE=eth0\nDNS=\n",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "empty data",
			iface:  "eth0",
			data:   "",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "real systemd-networkd lease excerpt",
			iface:  "eth0",
			data: `INTERFACE=eth0
DNS=1.1.1.1 1.0.0.1
DOMAINS=~
NTP=0.pool.ntp.org
`,
			want:   []string{"1.1.1.1", "1.0.0.1"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseSystemdNetworkdContent(tt.iface, []byte(tt.data))
			if ok != tt.wantOK {
				t.Errorf("parseSystemdNetworkdContent() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSystemdNetworkdContent() = %v, want %v", got, tt.want)
			}
		})
	}
}


