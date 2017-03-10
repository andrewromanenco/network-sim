package nsim

import (
	"net"
	"testing"
)

// func testARPNode() Node {
// 	node, _ := NewNodeBuilder().
// 		AddNetInterface("192.168.1.1/24").
// 		AddNetInterface("192.168.2.2/24").
// 		WithMedium(&dummyMedium{}).
// 		Build()
// 	node.AddRoute("192.168.3.0/24", "192.168.1.100")
// 	return node
// }

func TestARPReturnsNetworkInterfaceIPForSameNetwork(t *testing.T) {
	node := NewMockNode(t)
	// node.FRoutingTable = func() []Route {
	// 	return []Route{ParseRoute("192.168.3.0/24", "192.168.1.100")}
	// }
	ip := net.ParseIP("192.168.1.99")
	mac := ARP(node, ip)
	if mac != "192.168.1.99" {
		t.Error("MAC must be resolved to same IP as request. Because it's part of network interface network.")
	}
}

func TestARPReturnsDestinationForExistingRoute(t *testing.T) {
	node := NewMockNode(t)
	ip := net.ParseIP("192.168.3.33")
	mac := ARP(node, ip)
	if mac != "192.168.1.100" {
		t.Error("MAC must be resolved to IP of route destination.")
	}
}

func TestARPReturnsNothingForNonExistingRoute(t *testing.T) {
	node := NewMockNode(t)
	ip := net.ParseIP("192.168.5.55")
	mac := ARP(node, ip)
	if mac != "" {
		t.Error("Resolution must fail for non existing route.")
	}
}

func testARPPicksMostSpecificRoute(t *testing.T) {
	node := NewMockNode(t)
	//node.AddRoute("192.168.3.0/28", "192.168.1.200")
	ip := net.ParseIP("192.168.3.33")
	mac := ARP(node, ip)
	if mac != "192.168.1.200" {
		t.Error("Most specific route must be used.")
	}
}
