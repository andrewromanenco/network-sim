package nsim

import (
	"net"
	"testing"
)

type MockICMPPacket struct {
	MockIPPacket
	FMessageType func() string
}

func (micmp *MockICMPPacket) MessageType() string {
	if micmp.FMessageType == nil {
		panic("Message type handler is not set for ICMP mock packet.")
	}
	return micmp.FMessageType()
}

func TestICMPEcho(t *testing.T) {
	node := NewMockNode(t)
	replied := false
	var reply IPPacket
	fIPSend = func(n Node, p IPPacket) error {
		reply = p
		replied = true
		return nil
	}
	packet := MockICMPPacket{*NewMockIPPacket(t), nil}
	packet.FMessageType = func() string { return ICMPEchoRequest }
	packet.FSource = func() net.IP { return net.ParseIP("192.168.1.1") }
	packet.FDestination = func() net.IP { return net.ParseIP("192.168.2.2") }
	ICMPHandler(node, &packet)
	if !replied {
		t.Error("Nothing was sent back.")
	}
	icmpReply, ok := reply.(ICMPPacket)
	if !ok {
		t.Error("ICMP packet is expected as a reply.")
	}
	if !icmpReply.Destination().Equal(net.ParseIP("192.168.1.1")) {
		t.Error("Wrong reply destination.")
	}
	if !icmpReply.Source().Equal(net.ParseIP("192.168.2.2")) {
		t.Error("Wrong reply source.")
	}
	if icmpReply.MessageType() != ICMPEchoReply {
		t.Error("Echo reply is expected.")
	}
}

func TestNewICMPPacket(t *testing.T) {
	ipPacket := NewMockIPPacket(t)
	icmpPacket := NewICMPPacket(ipPacket, "msg-type")
	if icmpPacket == nil {
		t.Error("Must not be nil.")
		return
	}
	if icmpPacket.MessageType() != "msg-type" {
		t.Error("ICMP type is not correct.")
	}
	icmpStruct, ok := icmpPacket.(*icmp)
	if !ok {
		t.Error("icmp struct is expeced.")
		return
	}
	if icmpStruct.IPPacket != ipPacket {
		t.Error("IPPacket was not set.")
	}
}
