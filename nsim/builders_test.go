package nsim

import (
	"net"
	"testing"
)

func TestParseNetworkInterface(t *testing.T) {
	ni := ParseNetworkInterface("192.168.5.1/24")
	if ni == nil {
		t.Error("Valid CIDR must create a network interface.")
	}
	if !ni.IP.Equal(net.ParseIP("192.168.5.1")) {
		t.Error("Network interface does not have expected ip address.")
	}
	if !ni.Network.IP.Equal(net.ParseIP("192.168.5.0")) {
		t.Error("Network interface does not have expected network ip address.")
	}
	if size, _ := ni.Network.Mask.Size(); size != 24 {
		t.Error("Network interface does not have expected mask.")
	}
}

func TestParseNetworkInterfaceFails(t *testing.T) {
	if ParseNetworkInterface("192.168.5.1/200") != nil {
		t.Error("Expected to fail.")
	}
	if ParseNetworkInterface("not a cidr") != nil {
		t.Error("Expected to fail.")
	}
	if ParseNetworkInterface("192.168.5.1") != nil {
		t.Error("Expected to fail.")
	}
	if ParseNetworkInterface("/24") != nil {
		t.Error("Expected to fail.")
	}
}

func TestIPPacketBuilder(t *testing.T) {
	node := NewNodeBuilder().
		AddNetInterface("192.168.1.1/24").
		AddNetInterface("192.168.2.2/24").
		AddRoute("192.168.6.6/16", "192.168.6.1").
		AddRoute("192.168.5.6/16", "192.168.5.1").
		WithMedium(&dummyMedium{}).
		Build()
	if node == nil {
		t.Error("Result can not be null.")
	}
	if len(node.NetworkInterfaces()) != 2 {
		t.Error("Node must have two network interfaces.")
	}
	if !node.NetworkInterfaces()[0].Equal(ParseNetworkInterface("192.168.1.1/24")) {
		t.Error("Network interface config does not match.")
	}
	if len(node.RoutingTable()) != 2 {
		t.Error("Node must have two routes.")
	}
	if !node.RoutingTable()[0].Equal(ParseRoute("192.168.6.6/16", "192.168.6.1")) {
		t.Error("First route does not match config.")
	}
	if node.Medium() == nil {
		t.Error("Medium should not be null.")
	}
}
