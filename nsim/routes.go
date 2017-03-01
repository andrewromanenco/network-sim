package nsim

import (
	"errors"
	"net"
)

var (
	// ErrRouteInvalidCIDR means that provided network is not a valid one.
	ErrRouteInvalidCIDR = errors.New("Not a valid CIDR for a route.")

	// ErrRouteInvalidDestinationIP means that IP is not a valid one.
	ErrRouteInvalidDestinationIP = errors.New("Not a valid IP destination for a route.")
)

// AddRoute adds a route to node's routing table.
func (node *Node) AddRoute(cidrNet string, destinationIP string) error {
	ip := net.ParseIP(destinationIP)
	if ip == nil {
		return ErrRouteInvalidDestinationIP
	}
	_, network, err := net.ParseCIDR(cidrNet)
	if err != nil {
		return ErrRouteInvalidCIDR
	}
	node.RoutingTable = append(node.RoutingTable, Route{ip, *network})
	return nil
}
