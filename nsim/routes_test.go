package nsim

import "testing"

func TestAddRouteAddsItem(t *testing.T) {
	node := Node{nil, &dummyMedium{}, nil}
	err := node.AddRoute("192.168.1.0/24", "1.2.3.4")
	if err != nil {
		t.Error("Valid route should be added with no errors.")
	}
	if len(node.RoutingTable) != 1 {
		t.Error("Route was not added to routing table.")
	}
}

func TestAddRouteKeepsRoutesSortedByMask(t *testing.T) {
	node := Node{nil, &dummyMedium{}, nil}
	node.AddRoute("192.168.1.0/16", "2.2.2.2")
	node.AddRoute("192.168.1.0/8", "3.3.3.3")
	node.AddRoute("192.168.1.0/24", "1.1.1.1")
	if node.RoutingTable[0].DestinationIP.String() != "1.1.1.1" ||
		node.RoutingTable[1].DestinationIP.String() != "2.2.2.2" ||
		node.RoutingTable[2].DestinationIP.String() != "3.3.3.3" {
		t.Error("Routes must be sorted by mask size desc.")
	}
}

func TestAddRouteFailsIfCIDRisNotCorrect(t *testing.T) {
	node := Node{nil, &dummyMedium{}, nil}
	err := node.AddRoute("not-a-cidr", "1.2.3.4")
	if err != ErrRouteInvalidCIDR {
		t.Error("Valid route should be added with no errors.")
	}
}

func TestAddRouteFailsIfDestinationIsNotCorrectIP(t *testing.T) {
	node := Node{nil, &dummyMedium{}, nil}
	err := node.AddRoute("192.168.1.0/24", "not-an-ip")
	if err != ErrRouteInvalidDestinationIP {
		t.Error("Valid route should be added with no errors.")
	}
}
