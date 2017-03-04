package nsim

import "testing"

func TestDropFrameIfNotATarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.100.100", IPPacket{}}
	if node.LinkReceive(frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10/24").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.0.10", IPPacket{}}
	if !node.LinkReceive(frame) {
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
	if !node.LinkReceive(frame) {
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
	node.LinkSend(frame)
	acceptedFrame := *medium.acceptedFrame
	if acceptedFrame.destinationID != frame.destinationID {
		t.Error("Frame sent does not contained passed destination id.")
	}
	if acceptedFrame.IPPacket.TTL != frame.IPPacket.TTL {
		t.Error("Seems like IP packet was changed.")
	}
}
