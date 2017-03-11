package nsim

import (
	"net"
	"testing"
)

func testARPNode(t *testing.T) *MockNode {
	node := NewMockNode(t)
	node.FNetworkInterfaces = func() []NetworkInterface {
		return []NetworkInterface{
			*ParseNetworkInterface("192.168.1.1/24"),
			*ParseNetworkInterface("192.168.2.2/24"),
		}
	}
	node.FRoutingTable = func() []Route {
		return []Route{
			*ParseRoute("192.168.3.0/24", "192.168.1.100"),
		}
	}
	return node
}

func TestARPReturnsNetworkInterfaceIPForSameNetwork(t *testing.T) {
	node := testARPNode(t)
	ip := net.ParseIP("192.168.1.99")
	mac := ARP(node, ip)
	if mac != "192.168.1.99" {
		t.Error("MAC must be resolved to same IP as request. Because it's part of network interface network.")
	}
}

func TestARPReturnsDestinationForExistingRoute(t *testing.T) {
	node := testARPNode(t)
	ip := net.ParseIP("192.168.3.33")
	mac := ARP(node, ip)
	if mac != "192.168.1.100" {
		t.Error("MAC must be resolved to IP of route destination.")
	}
}

func TestARPReturnsNothingForNonExistingRoute(t *testing.T) {
	node := testARPNode(t)
	ip := net.ParseIP("192.168.5.55")
	mac := ARP(node, ip)
	if mac != "" {
		t.Error("Resolution must fail for non existing route.")
	}
}

func testARPChecksRoutesInOrder(t *testing.T) {
	node := testARPNode(t)
	node.FRoutingTable = func() []Route {
		return []Route{
			*ParseRoute("192.168.3.0/24", "192.168.1.200"),
			*ParseRoute("192.168.3.0/24", "192.168.1.100"),
		}
	}
	ip := net.ParseIP("192.168.3.33")
	mac := ARP(node, ip)
	if mac != "192.168.1.200" {
		t.Error("Most specific route must be used.")
	}
}
