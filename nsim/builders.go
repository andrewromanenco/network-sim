package nsim

import "net"

// ParseNetworkInterface creates a network interface from CIDR.
func ParseNetworkInterface(cidr string) *NetworkInterface {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &NetworkInterface{ip, *net}
}
