package nsim

import (
	"net"
	"testing"
)

func mockARP(node Node, ip net.IP) string {
	return "dest-mack"
}

func TestIPSendToExistingNode(t *testing.T) {
	node := NewMockNode(t)
	ipPacket := NewMockIPPacket(t)
	ipPacket.FDestination = func() net.IP { return net.ParseIP("192.168.0.1") }
	fARP = mockARP
	var sentFrame *Frame
	fLinkSend = func(node Node, frame Frame) error {
		sentFrame = &frame
		return nil
	}
	err := IPSend(node, ipPacket)
	if err != nil {
		t.Error("IPSend failed for no reason.")
	}
	if sentFrame == nil {
		t.Error("Seems like frame was not sent")
	}
}

func TestIPSendToNonExistingNodeFails(t *testing.T) {
	node := NewMockNode(t)
	ipPacket := NewMockIPPacket(t)
	ipPacket.FDestination = func() net.IP { return net.ParseIP("192.168.0.1") }
	fARP = func(node Node, ip net.IP) string {
		return ""
	}
	err := IPSend(node, ipPacket)
	if err != ErrNoRoute {
		t.Error("IPSend must fail if no next hop is found.")
	}
}

func TestIPSendPacketNoDestinationIPFails(t *testing.T) {
	node := NewMockNode(t)
	ipPacket := NewMockIPPacket(t)
	ipPacket.FDestination = func() net.IP { return nil }
	err := IPSend(node, ipPacket)
	if err != ErrIPDestinationNotSet {
		t.Error("IPSend must fail if packet has no destination ip.")
	}
}

func TestIPReceiveForwardsPacketIfRouter(t *testing.T) {
	var forwardedPacket IPPacket
	forwarded := false
	fIPSend = func(node Node, packet IPPacket) error {
		forwarded = true
		forwardedPacket = packet
		return nil
	}
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
			*ParseNetworkInterface("192.168.2.2/24"),
		}
	}
	packet := NewMockIPPacket(t)
	packet.FTTL = func() int { return 10 }
	packet.FDestination = func() net.IP { return net.ParseIP("192.168.10.10") }
	decreased := false
	packet.FDecreaseTTL = func() { decreased = true }
	IPReceive(node, packet)
	if !forwarded {
		t.Error("Packet was not forwarded.")
	}
	if !decreased {
		t.Error("Forwarded packet must have decreased TTL.")
	}
}

func TestIPReceiveNoForwardingIfNotARouter(t *testing.T) {
	forwarded := false
	fIPSend = func(node Node, packet IPPacket) error {
		forwarded = true
		return nil
	}
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
		}
	}
	packet := NewMockIPPacket(t)
	packet.FDestination = func() net.IP { return net.ParseIP("192.168.10.10") }
	IPReceive(node, packet)
	if forwarded {
		t.Error("Packet should not be forwarded if not a router.")
	}
}

func TestIPReceiveCallsProtocolHandlerWhenNodeIsDestination(t *testing.T) {
	forwarded := false
	fIPSend = func(node Node, packet IPPacket) error {
		forwarded = true
		return nil
	}
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
			*ParseNetworkInterface("192.168.2.2/24"),
		}
	}
	packet := NewMockIPPacket(t)
	packet.FDestination = func() net.IP { return net.ParseIP("192.168.1.1") }
	packet.FProtocol = func() string { return "some-protocol" }
	handled := false
	RegisterProtocolHandler("some-protocol", func(n Node, p IPPacket) { handled = true })
	IPReceive(node, packet)
	if forwarded {
		t.Error("Packet should not be forwarded, as it's destination is the node.")
	}
	if !handled {
		t.Error("Packet must be handled by protocol handler.")
	}
}

func TestRegisterProtocolHandler(t *testing.T) {
	count := 0
	var handler ProtocolHandler = func(n Node, p IPPacket) { count++ }
	RegisterProtocolHandler("some-protocol", handler)
	result := UnregisterProtocolHandler("some-protocol")
	result(nil, nil)
	if count != 1 {
		t.Error("Register must set handler, and unregister must return it.")
	}
}

func TestUnregisterProtocolHandler(t *testing.T) {
	if UnregisterProtocolHandler("nothing") != nil {
		t.Error("Must be nil for non existing handler.")
	}
}
