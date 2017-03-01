package nsim

import "testing"

func TestNodeBuilderFailsIfNoNetInterfaces(t *testing.T) {
	testee := NewNodeBuilder()
	_, err := testee.WithMedium(&dummyMedium{}).Build()
	if err != ErrNoNetworkInterfaces {
		t.Error("Error is expected when no interfaces are provided.")
	}
}

func TestNodeBuilderFailsIfNetInterfacesWithInvalidIP(t *testing.T) {
	testee := NewNodeBuilder()
	_, err := testee.AddNetInterface("not.an.ip").Build()
	if err != ErrNetworkInterfacesBadIP {
		t.Error("Error is expected when ip is not valid.")
	}
}

func TestNodeBuilderFailsIfNetInterfacesJustAnIP(t *testing.T) {
	testee := NewNodeBuilder()
	_, err := testee.AddNetInterface("192.168.0.1").Build()
	if err != ErrNetworkInterfacesBadIP {
		t.Error("Must fail if only IP is provided for an interface.")
	}
}

func TestNodeBuilderFailsIfNoMediumIsProvided(t *testing.T) {
	testee := NewNodeBuilder()
	_, err := testee.
		AddNetInterface("192.168.0.1/24").
		Build()
	if err != ErrNoTransmissionMedium {
		t.Error("Error is expected when no medium is provided.")
	}
}

func TestNodeBuilder(t *testing.T) {
	testee := NewNodeBuilder()
	node, err := testee.
		AddNetInterface("192.168.1.1/24").
		AddNetInterface("192.168.2.2/24").
		WithMedium(&dummyMedium{}).
		Build()
	if err != nil {
		t.Error("No error is expected for a good config.")
	}
	if len(node.NetworkInterfaces) != 2 {
		t.Error("Node must have two network interfaces")
	}
	if node.NetworkInterfaces[0].IP.String() != "192.168.1.1" {
		t.Error("First IP does not match the provided config.")
	}
	if node.NetworkInterfaces[0].Network.IP.String() != "192.168.1.0" {
		t.Error("First IP network does not match to the provided config.")
	}
	if size, _ := node.NetworkInterfaces[0].Network.Mask.Size(); size != 24 {
		t.Error("First Mask does not match to the provided config.")
	}
	if node.NetworkInterfaces[1].IP.String() != "192.168.2.2" {
		t.Error("Second IP does not match the provided config.")
	}
	if node.NetworkInterfaces[1].Network.IP.String() != "192.168.2.0" {
		t.Error("Second IP network does not match to the provided config.")
	}
	if size, _ := node.NetworkInterfaces[1].Network.Mask.Size(); size != 24 {
		t.Error("Second Mask does not match to the provided config.")
	}
	if node.Medium == nil {
		t.Error("Medium was not set.")
	}
}
