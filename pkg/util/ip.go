package util

import (
	"net"
)

// GetLocalIP returns a local IP address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() == nil {
			continue
		}

		return ipnet.IP.String()
	}

	return ""
}
