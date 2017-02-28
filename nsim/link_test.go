package nsim

import "testing"

func TestDropFrameIfNotATarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10").Build()
	frame := Frame{"192.168.100.100"}
	if node.LinkReceive(frame) {
		t.Error("Frame with other destination must be declined.")
	}
}

func TestAcceptFrameIfTarget(t *testing.T) {
	node, _ := NewNodeBuilder().AddNetInterface("192.168.0.10").Build()
	frame := Frame{"192.168.0.10"}
	if !node.LinkReceive(frame) {
		t.Error("Frame must be accepted.")
	}
}

func TestAcceptFrameIfTargetWithMultipleIPs(t *testing.T) {
	node, _ := NewNodeBuilder().
		AddNetInterface("192.168.0.10").
		AddNetInterface("192.168.1.10").
		Build()
	frame := Frame{"192.168.0.10"}
	if !node.LinkReceive(frame) {
		t.Error("Frame must be accepted.")
	}
}
