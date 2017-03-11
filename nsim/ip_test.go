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
	IPReceive(node, packet)
	if forwarded {
		t.Error("Packet should not be forwarded if not a router.")
	}
}
