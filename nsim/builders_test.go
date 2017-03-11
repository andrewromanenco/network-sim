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
