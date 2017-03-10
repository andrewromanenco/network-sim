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
