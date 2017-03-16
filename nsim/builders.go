package nsim

import "net"

// NodeBuilder is a builder for a node.
type NodeBuilder struct {
	networkInterfaces  []string
	medium             TransmissionMedium
	routes             []string
	routesDestinations []string
}

// NewNodeBuilder creates new builder.
func NewNodeBuilder() *NodeBuilder {
	return &NodeBuilder{}
}

// AddNetInterface adds network interface to this builder.
func (b *NodeBuilder) AddNetInterface(cidr string) *NodeBuilder {
	b.networkInterfaces = append(b.networkInterfaces, cidr)
	return b
}

// AddRoute adds route to this builder.
func (b *NodeBuilder) AddRoute(cidr string, gateway string) *NodeBuilder {
	b.routes = append(b.routes, cidr)
	b.routesDestinations = append(b.routesDestinations, gateway)
	return b
}

// WithMedium sets medium to be used on node creation.
func (b *NodeBuilder) WithMedium(medium TransmissionMedium) *NodeBuilder {
	b.medium = medium
	return b
}

// Build builds a node, or nil in case of any errors.
func (b *NodeBuilder) Build() Node {
	if b.medium == nil {
		return nil
	}
	ni := configureNetworkInterfaces(b.networkInterfaces)
	if ni == nil {
		return nil
	}
	routingTable := configureRoutingTable(b.routes, b.routesDestinations)
	if routingTable == nil {
		return nil
	}
	return &node{ni, b.medium, routingTable}
}

func configureNetworkInterfaces(configs []string) []NetworkInterface {
	if len(configs) == 0 {
		return nil
	}
	var ni []NetworkInterface
	for _, config := range configs {
		interf := ParseNetworkInterface(config)
		if interf == nil {
			return nil
		}
		ni = append(ni, *interf)
	}
	return ni
}

func configureRoutingTable(routes []string, destinations []string) []Route {
	if len(routes) == 0 {
		return nil
	}
	var routingTable []Route
	for index, r := range routes {
		route := ParseRoute(r, destinations[index])
		if route == nil {
			return nil
		}
		routingTable = append(routingTable, *route)
	}
	return routingTable
}

// ParseNetworkInterface creates a network interface from CIDR.
func ParseNetworkInterface(cidr string) *NetworkInterface {
	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &NetworkInterface{ip, *net}
}

// IPPacketBuilder is a builder of a IP packet.
type IPPacketBuilder struct {
	destination net.IP
	source      net.IP
	ttl         int
	protocol    string
}

// NewIPPacketBuilder returns a new ip packet builder.
func NewIPPacketBuilder() *IPPacketBuilder {
	return &IPPacketBuilder{}
}

// Destination sets destination for packet to be created.
func (ipb *IPPacketBuilder) Destination(ip string) *IPPacketBuilder {
	ipb.destination = net.ParseIP(ip)
	return ipb
}

// DestinationIP sets destination for packet to be created.
func (ipb *IPPacketBuilder) DestinationIP(ip net.IP) *IPPacketBuilder {
	ipb.destination = ip
	return ipb
}

// Source sets source for packet to be created.
func (ipb *IPPacketBuilder) Source(ip string) *IPPacketBuilder {
	ipb.source = net.ParseIP(ip)
	return ipb
}

// SourceIP sets source for packet to be created.
func (ipb *IPPacketBuilder) SourceIP(ip net.IP) *IPPacketBuilder {
	ipb.source = ip
	return ipb
}

// TTL sets ttl for packet to be created.
func (ipb *IPPacketBuilder) TTL(ttl int) *IPPacketBuilder {
	ipb.ttl = ttl
	return ipb
}

// Protocol sets protocol for packet to be created.
func (ipb *IPPacketBuilder) Protocol(protocol string) *IPPacketBuilder {
	ipb.protocol = protocol
	return ipb
}

// Build builds packet.
func (ipb *IPPacketBuilder) Build() IPPacket {
	return &ipPacket{ipb.destination, ipb.source, ipb.ttl, ipb.protocol}
}
