package nsim

import (
	"errors"
	"net"
	"reflect"
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

// Equals tests if two packets are equal.
func (packet *IPPacket) Equals(other *IPPacket) bool {
	if other == nil {
		return false
	}
	if packet.TTL != other.TTL {
		return false
	}
	if packet.Protocol != other.Protocol {
		return false
	}
	if !reflect.DeepEqual(packet.Destination, other.Destination) {
		return false
	}
	if !reflect.DeepEqual(packet.Source, other.Source) {
		return false
	}
	return true
}

var fARP = ARP
var fLinkSend = LinkSend

// IPSend sends an IP packet.
func IPSend(node *Node, packet IPPacket) error {
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

func isRouter(node *Node) bool {
	return len(node.NetworkInterfaces) > 1
}
