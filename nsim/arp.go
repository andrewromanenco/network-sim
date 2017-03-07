package nsim

import "net"

// ARPProto is a function type for an arp protocol implementaion.
type ARPProto func(node *Node, ip *net.IP)

// ARP resolves an IP to link layer identifier. The result is either a host in the
// same network as on of node's interfaces, or destination ID of a router. Or nothing.
// If router is not in the same network as the node, result is nothing.
func ARP(node *Node, ip *net.IP) string {
	netInterface := getNetworkInterfaceWithSameNetwork(node, ip)
	if netInterface != nil {
		return ip.String()
	}
	routerIP := getRouterIP(node, ip)
	if routerIP == nil {
		return ""
	}
	netInterface = getNetworkInterfaceWithSameNetwork(node, routerIP)
	if netInterface != nil {
		return routerIP.String()
	}
	return ""
}

func getNetworkInterfaceWithSameNetwork(node *Node, ip *net.IP) *NetworkInterface {
	for _, nInterface := range node.NetworkInterfaces {
		if nInterface.Network.Contains(*ip) {
			return &nInterface
		}
	}
	return nil
}

func getRouterIP(node *Node, ip *net.IP) *net.IP {
	for _, route := range node.RoutingTable {
		if route.Network.Contains(*ip) {
			return &route.DestinationIP
		}
	}
	return nil
}
