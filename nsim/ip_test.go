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
	err := IPSend(node, ipPacket)
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
	err := IPSend(node, ipPacket)
	if err != ErrNoRoute {
		t.Error("IPSend must fail if no next hop is found.")
	}
}

func TestIPSendPacketNoDestinationIPFails(t *testing.T) {
	node := testIPNode(t)
	ipPacket := testIPPacket()
	ipPacket.Destination = nil
	err := IPSend(node, ipPacket)
	if err != ErrIPDestinationNotSet {
		t.Error("IPSend must fail if packet has no destination ip.")
	}
}

func TestIPPacketEqualsWhenSame(t *testing.T) {
	packet1 := testIPPacket()
	packet2 := testIPPacket()
	if !packet1.Equals(&packet2) {
		t.Error("Equal packets must be equal.")
	}
}

func TestIPPacketEqualsForDefault(t *testing.T) {
	packet1 := IPPacket{}
	packet2 := IPPacket{}
	if !packet1.Equals(&packet2) {
		t.Error("Equal packets must be equal.")
	}
}

func TestIPPacketEqualsFalseForNil(t *testing.T) {
	packet1 := testIPPacket()
	if packet1.Equals(nil) {
		t.Error("Nil must be false.")
	}
}

func TestIPPacketNotEqualsDestination(t *testing.T) {
	packet1 := testIPPacket()
	packet2 := testIPPacket()
	packet2.Destination = net.ParseIP("127.0.0.1")
	if packet1.Equals(&packet2) {
		t.Error("Packets are different.")
	}
}

func TestIPPacketNotEqualSource(t *testing.T) {
	packet1 := testIPPacket()
	packet2 := testIPPacket()
	packet2.Source = net.ParseIP("127.0.0.1")
	if packet1.Equals(&packet2) {
		t.Error("Packets are different.")
	}
}

func TestIPPacketNotEqualTTL(t *testing.T) {
	packet1 := testIPPacket()
	packet2 := testIPPacket()
	packet2.TTL = 99
	if packet1.Equals(&packet2) {
		t.Error("Packets are different.")
	}
}

func TestIPPacketNotEqualProtocol(t *testing.T) {
	packet1 := testIPPacket()
	packet2 := testIPPacket()
	packet2.Protocol = "other"
	if packet1.Equals(&packet2) {
		t.Error("Packets are different.")
	}
}

func TestIsNotARouterIfOneInterfaceOnly(t *testing.T) {
	node := testIPNode(t)
	node.NetworkInterfaces = node.NetworkInterfaces[:len(node.NetworkInterfaces)-1]
	if isRouter(node) {
		t.Error("Node with one net interface can not be a router")
	}
}

func TestIsARouterIfSeveralInterfaces(t *testing.T) {
	node := testIPNode(t)
	if !isRouter(node) {
		t.Error("Node with multiple net interfaces mut be a router")
	}
}

func TestNodeOwnsIP(t *testing.T) {
	node := testIPNode(t)
	ip := net.ParseIP("192.168.1.1")
	if !nodeOwnsIP(node, &ip) {
		t.Error("Must return true as node has this ip.")
	}
}

func TestNodeNotOwnsIP(t *testing.T) {
	node := testIPNode(t)
	ip := net.ParseIP("192.168.100.100")
	if nodeOwnsIP(node, &ip) {
		t.Error("Must return false for unknown ip.")
	}
}
