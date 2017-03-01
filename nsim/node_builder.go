package nsim

import (
	"errors"
	"net"
)

var (
	//ErrNoNetworkInterfaces means that no network interface configs were provided.
	ErrNoNetworkInterfaces = errors.New("No network interfaces provided. At least one is required.")

	//ErrNetworkInterfacesBadIP means that at least on IP config for interfaces is not a valid one.
	ErrNetworkInterfacesBadIP = errors.New("Provide IP/Mask as CIDR.")

	//ErrNoTransmissionMedium means that medium was not provided.
	ErrNoTransmissionMedium = errors.New("No transmission medium is provided.")
)

// NodeBuilder builds a node.
type NodeBuilder struct {
	networkInterfaces []string
	medium            TransmissionMedium
}

// NewNodeBuilder creates new NodeBuilder.
func NewNodeBuilder() *NodeBuilder {
	return &NodeBuilder{make([]string, 0), nil}
}

// AddNetInterface adds network interface to the node under construction.
func (nb *NodeBuilder) AddNetInterface(ip string) *NodeBuilder {
	nb.networkInterfaces = append(nb.networkInterfaces, ip)
	return nb
}

// WithMedium sets medium to be used for data exchange.
func (nb *NodeBuilder) WithMedium(medium TransmissionMedium) *NodeBuilder {
	nb.medium = medium
	return nb
}

// Build creates a valid node. Or returns an error.
func (nb *NodeBuilder) Build() (*Node, error) {
	if len(nb.networkInterfaces) == 0 {
		return nil, ErrNoNetworkInterfaces
	}
	var networkInterfaces []NetworkInterface
	for _, ip := range nb.networkInterfaces {
		interfaceIP, interfaceNet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, ErrNetworkInterfacesBadIP
		}
		networkInterfaces = append(networkInterfaces, NetworkInterface{interfaceIP, *interfaceNet})
	}
	if nb.medium == nil {
		return nil, ErrNoTransmissionMedium
	}
	return &Node{networkInterfaces, nb.medium}, nil
}
