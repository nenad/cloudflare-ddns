package ip_test

import (
	"cloudflare-ddns/pkg/ip"
	"net"
	"strings"
	"testing"

	"golang.org/x/net/nettest"
)

func TestGetAddress(t *testing.T) {
	tests := []struct {
		name      string
		version   ip.Version
		flags     net.Flags
		isUnicast bool
		err       string
	}{
		{
			name:      "ipv6 interface up will return unicast address",
			version:   ip.V6,
			flags:     net.FlagUp,
			isUnicast: true,
		},
		{
			name:      "ipv4 interface up will return unicast address",
			version:   ip.V4,
			flags:     net.FlagMulticast,
			isUnicast: true,
		},
		{
			name:      "ip loopback interface up will return no addresses",
			version:   ip.V4,
			flags:     net.FlagUp | net.FlagLoopback,
			isUnicast: false,
			err:       "could not find global unicast address for interface",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.version != "ip6" && !nettest.SupportsIPv4() {
				t.Skip("skipping test, IPv4 is not supported")
			}
			if tt.version == "ip6" && !nettest.SupportsIPv6() {
				t.Skip("skipping test, IPv6 is not supported")
			}

			iface, err := nettest.RoutedInterface(string(tt.version), tt.flags)
			if err != nil {
				t.Fatalf("could not create interface: %s", err)
			}

			retriever := ip.InterfaceRetriever{Device: iface.Name}

			ipAddr, err := retriever.Get(tt.version)
			if err != nil {
				if tt.err != "" {
					if !strings.Contains(err.Error(), tt.err) {
						t.Fatalf("expected error did not match: want %q, got %q", tt.err, err.Error())
					}
				} else {
					t.Fatalf("could not get IP address from interface: %s", err)
				}
			}

			addr := net.ParseIP(ipAddr)
			if addr.IsGlobalUnicast() != tt.isUnicast {
				t.Fatalf("device %q: expected unicast %t, got %t", iface.Name, tt.isUnicast, addr.IsGlobalUnicast())
			}
		})
	}
}
