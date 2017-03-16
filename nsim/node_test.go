package nsim

import "testing"

func TestNetworkInterfaceEqual(t *testing.T) {
	ni := ParseNetworkInterface("192.168.1.1/24")
	if ni.Equal(nil) {
		t.Error("Must not be equal")
	}
	if ni.Equal(ParseNetworkInterface("192.168.1.2/24")) {
		t.Error("Must not be equal")
	}
	if ni.Equal(ParseNetworkInterface("192.168.1.1/16")) {
		t.Error("Must not be equal")
	}
	if ni.Equal(ParseNetworkInterface("192.168.5.1/24")) {
		t.Error("Must not be equal")
	}
	if !ni.Equal(ParseNetworkInterface("192.168.1.1/24")) {
		t.Error("Must be equal")
	}
}
