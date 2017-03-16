package nsim

import "testing"

func TestHasMore(t *testing.T) {
	testee := NewSignleQueueMedium()
	if testee.HasMore() {
		t.Error("New medium may not have anything to process.")
	}
	testee.Send(Frame{"mac", nil})
	if !testee.HasMore() {
		t.Error("Must be true.")
	}
	testee.Send(Frame{"mac", nil})
	if !testee.HasMore() {
		t.Error("Must be true.")
	}
}

func TestProcessFrame(t *testing.T) {
	testee := NewSignleQueueMedium()
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
		}
	}
	testee.RegisterNode(node)
	frame := Frame{"192.168.1.1", nil}
	frameReived := false
	fLinkReceive = func(n Node, f Frame) bool {
		frameReived = true
		return true
	}
	testee.Send(frame)
	if testee.DeliverFrame() {
		t.Error("Should be false as only one frame exists.")
	}
	if !frameReived {
		t.Error("Frame was not delivered.")
	}
	testee.Send(frame)
	testee.Send(frame)
	if !testee.DeliverFrame() {
		t.Error("Should be true as one more frame exists.")
	}
	if testee.DeliverFrame() {
		t.Error("Should be false as no more frames to process.")
	}
}
