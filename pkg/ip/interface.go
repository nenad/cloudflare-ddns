package ip

import (
	"fmt"
	"net"
)

type InterfaceRetriever struct {
	Device string
}

// GetAddress returns the first unicast address for the given interface and protocol.
func (r *InterfaceRetriever) Get(version Version) (string, error) {
	iface, err := net.InterfaceByName(r.Device)
	if err != nil {
		return "", fmt.Errorf("could not get details about interface %q: %w", r.Device, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("could not get interface %q addresses: %w", r.Device, err)
	}

	for _, a := range addrs {
		switch v := a.(type) {
		case *net.IPNet:
			if version == "ip6" && v.IP.To4() != nil {
				continue
			}
			if v.IP.IsGlobalUnicast() {
				return v.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("could not find global unicast address for interface %q", r.Device)
}
