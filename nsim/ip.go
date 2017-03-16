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
	if packetForNode(node, packet) {
		if handler, ok := protocolHandlers[packet.Protocol()]; ok {
			handler(node, packet)
		}
		return
	}
	if !isRouter(node) {
		return
	}
	packet.DecreaseTTL()
	fIPSend(node, packet)
}

func isRouter(node Node) bool {
	return len(node.NetworkInterfaces()) > 1
}

func packetForNode(node Node, packet IPPacket) bool {
	for _, ni := range node.NetworkInterfaces() {
		if ni.IP.Equal(packet.Destination()) {
			return true
		}
	}
	return false
}

// ProtocolHandler is a protocl handler on top of IP layer.
type ProtocolHandler func(Node, IPPacket)

var protocolHandlers = make(map[string]ProtocolHandler)

// RegisterProtocolHandler registers a protocol handler on top of IP layer.
func RegisterProtocolHandler(protocol string, handler ProtocolHandler) {
	protocolHandlers[protocol] = handler
}

// UnregisterProtocolHandler unregisters a protocol handler on top of IP layer.
func UnregisterProtocolHandler(protocol string) ProtocolHandler {
	if handler, ok := protocolHandlers[protocol]; ok {
		delete(protocolHandlers, protocol)
		return handler
	}
	return nil
}
