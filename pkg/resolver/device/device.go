package device

import (
	"fmt"
	"net"
)

// GetAddress returns the first unicast address for the given interface and protocol.
func GetAddress(device, protocol string) (string, error) {
	iface, err := net.InterfaceByName(device)
	if err != nil {
		return "", fmt.Errorf("could not get details about interface %q: %w", device, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("could not get interface %q addresses: %w", device, err)
	}

	for _, a := range addrs {
		switch v := a.(type) {
		case *net.IPNet:
			if protocol == "ip6" && v.IP.To4() != nil {
				continue
			}
			if v.IP.IsGlobalUnicast() {
				return v.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("could not find global unicast address for interface %q", device)
}
