package nsim

import "net"

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
