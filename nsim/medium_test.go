package nsim

import "testing"

func TestRegisterNode(t *testing.T) {
	testee := NewSignleQueueMedium()
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
			*ParseNetworkInterface("192.168.2.2/24"),
		}
	}
	testee.RegisterNode(node)
}

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

func TestGetFrame(t *testing.T) {
	testee := NewSignleQueueMedium()
	if testee.Frame() != nil {
		t.Error("Must be nil as no frames available.")
	}
	testee.Send(Frame{"mac1", nil})
	if testee.Frame() == nil {
		t.Error("Must not be nil.")
	}
	testee.Send(Frame{"mac2", nil})
	if testee.Frame().destinationID != "mac1" {
		t.Error("First unprocessed frame must be returned.")
	}
	testee.DeliverFrame()
	if testee.Frame().destinationID != "mac2" {
		t.Error("Second must be available after first is processed.")
	}
	testee.DeliverFrame()
	if testee.Frame() != nil {
		t.Error("Should be none.")
	}
}
