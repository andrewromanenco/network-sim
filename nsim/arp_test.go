package nsim

import (
	"net"
	"testing"
)

func testARPNode() *Node {
	node, _ := NewNodeBuilder().
		AddNetInterface("192.168.1.1/24").
		AddNetInterface("192.168.2.2/24").
		WithMedium(&dummyMedium{}).
		Build()
	node.AddRoute("192.168.3.0/24", "192.168.1.100")
	return node
}

func TestARPReturnsNetworkInterfaceIPForSameNetwork(t *testing.T) {
	node := testARPNode()
	mac := node.ARP(net.ParseIP("192.168.1.99"))
	if mac != "192.168.1.99" {
		t.Error("MAC must be resolved to same IP as request. Because it's part of network interface network.")
	}
}

func TestARPReturnsDestinationForExistingRoute(t *testing.T) {
	node := testARPNode()
	mac := node.ARP(net.ParseIP("192.168.3.33"))
	if mac != "192.168.1.100" {
		t.Error("MAC must be resolved to IP of route destination.")
	}
}

func TestARPReturnsNothingForNonExistingRoute(t *testing.T) {
	node := testARPNode()
	mac := node.ARP(net.ParseIP("192.168.5.55"))
	if mac != "" {
		t.Error("Resolution must fail for non existing route.")
	}
}
