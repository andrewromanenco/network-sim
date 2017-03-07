package nsim

import (
	"errors"
	"net"
)

var (
	//ErrIPDestinationNotSet means that packet has no destination IP info.
	ErrIPDestinationNotSet = errors.New("No destination IP.")
)

// IPPacket is a packet routed over an IP network.
type IPPacket struct {
	Destination net.IP
	Source      net.IP
	TTL         int
	Protocol    string
}

var fARP = ARP

// IPSend sends an IP packet.
func (node *Node) IPSend(packet IPPacket) error {
	destinationMAC := fARP(node, &packet.Destination)
	frame := Frame{destinationMAC, packet}
	LinkSend(node, frame)
	return nil
}
