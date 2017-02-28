package nsim

import (
	"errors"
	"net"
)

var (
	//ErrNoNetworkInterfaces means that no network interface configs were provided.
	ErrNoNetworkInterfaces = errors.New("No network interfaces provided. At least one is required.")

	//ErrNetworkInterfacesBadIP means that at least on IP config for interfaces is not a valid one.
	ErrNetworkInterfacesBadIP = errors.New("Invalid IP address for an interface is provided.")
)

// NodeBuilder builds a node.
type NodeBuilder struct {
	networkInterfaces []string
}

// NewNodeBuilder creates new NodeBuilder.
func NewNodeBuilder() *NodeBuilder {
	return &NodeBuilder{make([]string, 0)}
}

// AddNetInterface adds network interface to the node under construction.
func (nb *NodeBuilder) AddNetInterface(ip string) *NodeBuilder {
	nb.networkInterfaces = append(nb.networkInterfaces, ip)
	return nb
}

// Build creates a valid node. Or returns an error.
func (nb *NodeBuilder) Build() (*Node, error) {
	if len(nb.networkInterfaces) == 0 {
		return nil, ErrNoNetworkInterfaces
	}
	var networkInterfaces []NetworkInterface
	for _, ip := range nb.networkInterfaces {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return nil, ErrNetworkInterfacesBadIP
		}
		networkInterfaces = append(networkInterfaces, NetworkInterface{parsedIP})
	}
	return &Node{networkInterfaces}, nil
}
