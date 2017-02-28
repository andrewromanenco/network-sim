package nsim

import (
	"net"
)

// Node is a model of a computer in the Internet.
type Node struct {
	NetworkInterfaces []NetworkInterface
}

// NetworkInterface is a model for a network card.
type NetworkInterface struct {
	IP net.IP
}
