package nsim

import (
	"errors"
	"net"
)

var (
	//ErrIPDestinationNotSet means that packet has no destination IP info.
	ErrIPDestinationNotSet = errors.New("No destination IP in the packet.")

	//ErrNoRoute means that there is no route to the destination.
	ErrNoRoute = errors.New("No route to destinaton.")
)

// IPPacket is a packet routed over an IP network.
type IPPacket struct {
	Destination net.IP
	Source      net.IP
	TTL         int
	Protocol    string
}

var fARP = ARP
var fLinkSend = LinkSend

// IPSend sends an IP packet.
func (node *Node) IPSend(packet IPPacket) error {
	if packet.Destination == nil {
		return ErrIPDestinationNotSet
	}
	destinationMAC := fARP(node, &packet.Destination)
	if destinationMAC == "" {
		return ErrNoRoute
	}
	frame := Frame{destinationMAC, packet}
	fLinkSend(node, frame)
	return nil
}
