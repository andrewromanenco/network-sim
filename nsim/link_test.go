package nsim

import "testing"

func TestDropFrameIfNotATarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.100.100"}
	if node.LinkReceive(frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10").WithMedium(&dummyMedium{}).Build()
	frame := Frame{"192.168.0.10"}
	if !node.LinkReceive(frame) {
		t.Error("Frame must be accepted.")
	}
}

func TestAcceptFrameIfTargetWithMultipleIPs(t *testing.T) {
	node, _ := NewNodeBuilder().
		AddNetInterface("192.168.0.10").
		AddNetInterface("192.168.1.10").
		WithMedium(&dummyMedium{}).
		Build()
	frame := Frame{"192.168.0.10"}
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
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10").WithMedium(medium).Build()
	frame := Frame{"192.168.0.10"}
	node.LinkSend(frame)
	if *medium.acceptedFrame != frame {
		t.Error("Frame was not sent to the medium.")
	}
}
