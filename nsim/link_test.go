package nsim

import (
	"net"
	"testing"
)

func init() {
	fIPReceive = func(node *Node, ipPacket IPPacket) {}
}

func TestDropFrameIfNotATarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.100.100", IPPacket{}}
	if LinkReceive(node, frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.0.10", IPPacket{}}
	if !LinkReceive(node, frame) {
		t.Error("Frame must be accepted.")
	}
}

func TestAcceptFrameIfTargetWithMultipleIPs(t *testing.T) {
	node, _ := NewNodeBuilder().
		AddNetInterface("192.168.0.10/24").
		AddNetInterface("192.168.1.10/24").
		WithMedium(&dummyMedium{}).
		Build()
	frame := Frame{"192.168.0.10", IPPacket{}}
	if !LinkReceive(node, frame) {
		t.Error("Frame must be accepted.")
	}
}

type mockMedium struct {
	acceptedFrame *Frame
}

func (m *mockMedium) send(frame Frame) error {
	m.acceptedFrame = &frame
	return nil
}

func TestLinkSend(t *testing.T) {
	medium := &mockMedium{nil}
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(medium).Build()
	frame := Frame{"192.168.0.10", IPPacket{}}
	frame.IPPacket.TTL = 999
	LinkSend(node, frame)
	acceptedFrame := *medium.acceptedFrame
	if acceptedFrame.destinationID != frame.destinationID {
		t.Error("Frame sent does not contained passed destination id.")
	}
	if acceptedFrame.IPPacket.TTL != frame.IPPacket.TTL {
		t.Error("Seems like IP packet was changed.")
	}
}

func TestFrameEqualsForSame(t *testing.T) {
	frame1 := Frame{"192.168.0.10", IPPacket{}}
	frame2 := Frame{"192.168.0.10", IPPacket{}}
	if !frame1.Equals(&frame2) {
		t.Error("Expected to be same.")
	}
}

func TestFrameEqualsFalseForNil(t *testing.T) {
	frame := Frame{"192.168.0.10", IPPacket{}}
	if frame.Equals(nil) {
		t.Error("Nil must be false.")
	}
}

func TestFrameEqualsFalseForDiffDestination(t *testing.T) {
	frame1 := Frame{"192.168.0.10", IPPacket{}}
	frame2 := Frame{"192.168.0.11", IPPacket{}}
	if frame1.Equals(&frame2) {
		t.Error("Expected to be different.")
	}
}

func TestFrameEqualsFalseForDiffPacket(t *testing.T) {
	frame1 := Frame{"192.168.0.10", IPPacket{}}
	frame2 := Frame{"192.168.0.10", IPPacket{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		100,
		"None",
	}}
	if frame1.Equals(&frame2) {
		t.Error("Expected to be different.")
	}
}
