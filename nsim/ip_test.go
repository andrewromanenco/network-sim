package nsim

import (
	"net"
	"testing"
)

func testIPNode(t *testing.T) *Node {
	node, err := NewNodeBuilder().
		AddNetInterface("192.168.1.1/24").
		AddNetInterface("192.168.2.2/24").
		WithMedium(&dummyMedium{}).
		Build()
	if err != nil {
		t.Error("Test Node for IP was not created.")
	}
	node.AddRoute("192.168.3.0/24", "192.168.1.100")
	return node
}

func testIPPacket() IPPacket {
	return IPPacket{
		net.ParseIP("192.168.0.1"),
		net.ParseIP("192.168.0.2"),
		10,
		"None",
	}
}

func mockARP(node *Node, ip *net.IP) string {
	return "dest-mack"
}

func TestIPSendToExistingNode(t *testing.T) {
	node := testIPNode(t)
	ipPacket := testIPPacket()
	fARP = mockARP
	var sentFrame *Frame
	fLinkSend = func(node *Node, frame Frame) error {
		sentFrame = &frame
		return nil
	}
	err := node.IPSend(ipPacket)
	if err != nil {
		t.Error("IPSend failed for no reason.")
	}
	if sentFrame == nil {
		t.Error("Seems like frame was not sent")
	}
}

func TestIPSendToNonExistingNodeFails(t *testing.T) {
	node := testIPNode(t)
	ipPacket := testIPPacket()
	fARP = func(node *Node, ip *net.IP) string {
		return ""
	}
	err := node.IPSend(ipPacket)
	if err != ErrNoRoute {
		t.Error("IPSend must fail if no next hop is found.")
	}
}

func TestIPSendPacketNoDestinationIPFails(t *testing.T) {
	node := testIPNode(t)
	ipPacket := testIPPacket()
	ipPacket.Destination = nil
	err := node.IPSend(ipPacket)
	if err != ErrIPDestinationNotSet {
		t.Error("IPSend must fail if packet has no destination ip.")
	}
}
