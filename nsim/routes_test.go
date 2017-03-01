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
