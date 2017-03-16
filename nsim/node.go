package nsim

import (
	"net"
)

// Node represent a computer in a network.
type Node interface {
	NetworkInterfaces() []NetworkInterface
	Medium() TransmissionMedium
	RoutingTable() []Route
}

type node struct {
	networkInterfaces []NetworkInterface
	medium            TransmissionMedium
	routingTable      []Route
}

func (n *node) NetworkInterfaces() []NetworkInterface {
	return n.networkInterfaces
}

func (n *node) Medium() TransmissionMedium {
	return n.medium
}

func (n *node) RoutingTable() []Route {
	return n.routingTable
}

// NetworkInterface is a model for a network card.
type NetworkInterface struct {
	IP      net.IP
	Network net.IPNet
}

// Route is a route to a specific network.
type Route struct {
	DestinationIP net.IP
	Network       net.IPNet
}

// Equal checks if two net interfaces are equal.
func (ni *NetworkInterface) Equal(other *NetworkInterface) bool {
	if other == nil {
		return false
	}
	if !ni.IP.Equal(other.IP) {
		return false
	}
	if !ni.Network.IP.Equal(other.Network.IP) {
		return false
	}
	thisSize, _ := ni.Network.Mask.Size()
	thatSize, _ := other.Network.Mask.Size()
	if thisSize != thatSize {
		return false
	}
	return true
}

// Equal checks if two routes are equal.
func (r *Route) Equal(other *Route) bool {
	if other == nil {
		return false
	}
	if !r.DestinationIP.Equal(other.DestinationIP) {
		return false
	}
	if !r.Network.IP.Equal(other.Network.IP) {
		return false
	}
	thisSize, _ := r.Network.Mask.Size()
	thatSize, _ := other.Network.Mask.Size()
	if thisSize != thatSize {
		return false
	}
	return true
}
