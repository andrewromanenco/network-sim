package nsim

import (
	"testing"
)

func init() {
	fIPReceive = func(node Node, ipPacket IPPacket) {}
}

func TestDropFrameIfNotATarget(t *testing.T) {
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface { return []NetworkInterface{*ParseNetworkInterface("192.168.0.10/24")} }
	frame := Frame{"192.168.100.100", NewMockIPPacket(t)}
	if LinkReceive(node, frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface { return []NetworkInterface{*ParseNetworkInterface("192.168.0.10/24")} }
	frame := Frame{"192.168.0.10", NewMockIPPacket(t)}
	if !LinkReceive(node, frame) {
		t.Error("Frame must be accepted.")
	}
}

func TestAcceptFrameIfTargetWithMultipleIPs(t *testing.T) {
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.0.10/24"),
			*ParseNetworkInterface("192.168.1.10/24"),
		}
	}
	frame := Frame{"192.168.0.10", NewMockIPPacket(t)}
	if !LinkReceive(node, frame) {
		t.Error("Frame must be accepted.")
	}
}

type mockMedium struct {
	acceptedFrame *Frame
}

func (m *mockMedium) Send(frame Frame) error {
	m.acceptedFrame = &frame
	return nil
}

func TestLinkSend(t *testing.T) {
	medium := &mockMedium{nil}
	node := NewMockNode(t)
	node.FMedium = func() TransmissionMedium { return medium }
	mPacket := NewMockIPPacket(t)
	mPacket.FTTL = func() int { return 999 }
	frame := Frame{"192.168.0.10", mPacket}
	LinkSend(node, frame)
	acceptedFrame := *medium.acceptedFrame
	if acceptedFrame.destinationID != frame.destinationID {
		t.Error("Frame sent does not contained passed destination id.")
	}
	if acceptedFrame.IPPacket.TTL() != frame.IPPacket.TTL() {
		t.Error("Seems like IP packet was changed.")
	}
}
