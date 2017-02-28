package nsim

import "testing"

func TestNodeBuilderFailsIfNoNetInterfaces(t *testing.T) {
	testee := NewNodeBuilder()
	_, err := testee.Build()
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

func TestNodeBuilder(t *testing.T) {
	testee := NewNodeBuilder()
	node, err := testee.
		AddNetInterface("192.168.0.1").
		AddNetInterface("192.168.0.2").
		Build()
	if err != nil {
		t.Error("No error is expected for a good config.")
	}
	if len(node.NetworkInterfaces) != 2 {
		t.Error("Node must have two network interfaces")
	}
	if node.NetworkInterfaces[0].IP.String() != "192.168.0.1" {
		t.Error("First IP does not match the provided config.")
	}
	if node.NetworkInterfaces[1].IP.String() != "192.168.0.2" {
		t.Error("Second IP does not match the provided config.")
	}
}
