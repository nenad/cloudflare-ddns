package ip

import (
	"time"
)

const (
	V4 Version = "ip4"
	V6 Version = "ip6"
)

type (
	Version   string
	Retriever interface {
		Get(version Version) (string, error)
	}
)

func Factory(iface string) Retriever {
	if iface != "" {
		return &InterfaceRetriever{Device: iface}
	}

	return NewClient(Retry(3), Timeout(time.Second*10))
}
