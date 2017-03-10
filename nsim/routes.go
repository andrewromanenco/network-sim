package nsim

import (
	"net"
)

// ParseRoute creates a route structure for a network.
func ParseRoute(cidrNet string, destinationIP string) *Route {
	ip := net.ParseIP(destinationIP)
	if ip == nil {
		return nil
	}
	_, network, err := net.ParseCIDR(cidrNet)
	if err != nil {
		return nil
	}
	return &Route{ip, *network}
}
