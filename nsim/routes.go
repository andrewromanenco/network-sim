package nsim

import (
	"errors"
	"net"
	"sort"
)

var (
	// ErrRouteInvalidCIDR means that provided network is not a valid one.
	ErrRouteInvalidCIDR = errors.New("Not a valid CIDR for a route.")

	// ErrRouteInvalidDestinationIP means that IP is not a valid one.
	ErrRouteInvalidDestinationIP = errors.New("Not a valid IP destination for a route.")
)

// AddRoute adds a route to node's routing table.
func (node *Node) AddRoute(cidrNet string, destinationIP string) error {
	ip := net.ParseIP(destinationIP)
	if ip == nil {
		return ErrRouteInvalidDestinationIP
	}
	_, network, err := net.ParseCIDR(cidrNet)
	if err != nil {
		return ErrRouteInvalidCIDR
	}
	node.RoutingTable = append(node.RoutingTable, Route{ip, *network})
	sortRoutesByMaskDesc(node.RoutingTable)
	return nil
}

type byMask []Route

func (rt byMask) Len() int      { return len(rt) }
func (rt byMask) Swap(i, j int) { rt[i], rt[j] = rt[j], rt[i] }
func (rt byMask) Less(i, j int) bool {
	iSize, _ := rt[i].Network.Mask.Size()
	jSize, _ := rt[j].Network.Mask.Size()
	return iSize > jSize
}

func sortRoutesByMaskDesc(routes []Route) {
	sort.Sort(byMask(routes))
}
