package nsim

import (
	"testing"
)

func init() {
	fIPReceive = func(node Node, ipPacket IPPacket) {}
}

func TestDropFrameIfNotATarget(t *testing.T) {
	//node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	node := NewMockNode(t)
	frame := Frame{"192.168.100.100", NewMockIPPacket(t)}
	if LinkReceive(node, frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	//node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	node := NewMockNode(t)
	frame := Frame{"192.168.0.10", NewMockIPPacket(t)}
	if !LinkReceive(node, frame) {
		t.Error("Frame must be accepted.")
	}
}

func TestAcceptFrameIfTargetWithMultipleIPs(t *testing.T) {
	// node, _ := NewNodeBuilder().
	// 	AddNetInterface("192.168.0.10/24").
	// 	AddNetInterface("192.168.1.10/24").
	// 	WithMedium(&dummyMedium{}).
	// 	Build()
	node := NewMockNode(t)
	frame := Frame{"192.168.0.10", NewMockIPPacket(t)}
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
	// node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(medium).Build()
	// frame := Frame{"192.168.0.10", IPPacket{}}
	node := NewMockNode(t)
	frame := Frame{"192.168.0.10", NewMockIPPacket(t)}
	//frame.IPPacket.TTL = 999
	LinkSend(node, frame)
	acceptedFrame := *medium.acceptedFrame
	if acceptedFrame.destinationID != frame.destinationID {
		t.Error("Frame sent does not contained passed destination id.")
	}
	if acceptedFrame.IPPacket.TTL() != frame.IPPacket.TTL() {
		t.Error("Seems like IP packet was changed.")
	}
}
