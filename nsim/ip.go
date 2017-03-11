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

// IPPacket represent a packet on IP layer of TCP/IP stack.
type IPPacket interface {
	Destination() net.IP
	Source() net.IP
	TTL() int
	Protocol() string
	DecreaseTTL()
}

type ipPacket struct {
	destination net.IP
	source      net.IP
	ttl         int
	protocol    string
}

func (ipp *ipPacket) Destination() net.IP {
	return ipp.destination
}

func (ipp *ipPacket) Source() net.IP {
	return ipp.source
}

func (ipp *ipPacket) TTL() int {
	return ipp.ttl
}

func (ipp *ipPacket) Protocol() string {
	return ipp.protocol
}

func (ipp *ipPacket) DecreaseTTL() {
	ipp.ttl--
}

var fARP = ARP
var fLinkSend = LinkSend
var fIPSend = IPSend

// IPSend sends an IP packet.
func IPSend(node Node, packet IPPacket) error {
	if packet.Destination() == nil {
		return ErrIPDestinationNotSet
	}
	destinationMAC := fARP(node, packet.Destination())
	if destinationMAC == "" {
		return ErrNoRoute
	}
	frame := Frame{destinationMAC, packet}
	fLinkSend(node, frame)
	return nil
}

// IPReceive is called when an IP packet arrives from lower layer.
func IPReceive(node Node, packet IPPacket) {
	if len(node.NetworkInterfaces()) == 1 {
		return
	}
	packet.DecreaseTTL()
	fIPSend(node, packet)
}
