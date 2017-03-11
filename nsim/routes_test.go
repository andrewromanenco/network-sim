package nsim

import (
	"net"
	"testing"
)

func TestParseRoute(t *testing.T) {
	route := ParseRoute("192.168.0.0/16", "192.168.0.1")
	if !route.DestinationIP.Equal(net.ParseIP("192.168.0.1")) {
		t.Error("Route doesn't have extected gateway.")
	}
	if !route.Network.IP.Equal(net.ParseIP("192.168.0.0")) {
		t.Error("Route doesn't have extected network ip.")
	}
	if size, _ := route.Network.Mask.Size(); size != 16 {
		t.Error("Route doesn't have extected network mask.")
	}
}
