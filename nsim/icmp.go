package nsim

const (
	// ICMPProtocol is constant id for ICMP packets over IP layer.
	ICMPProtocol = "ICMP"
	// ICMPEchoRequest represents ping request.
	ICMPEchoRequest = "echo-request"
	// ICMPEchoReply represents ping response.
	ICMPEchoReply = "echo-reply"
)

// ICMPPacket represents ICMP protocol.
type ICMPPacket interface {
	IPPacket
	MessageType() string
}

type icmp struct {
	IPPacket
	messageType string
}

// MessageType returns message type for an icmp packet.
func (ic *icmp) MessageType() string {
	return ic.messageType
}

// ICMPHandler handles incoming ICMP packets.
func ICMPHandler(node Node, packet IPPacket) {
	icmpPacket, ok := packet.(ICMPPacket)
	if !ok {
		panic("ICMP handler received non icmp packet")
	}
	if icmpPacket.MessageType() == ICMPEchoRequest {
		ipEcho := NewIPPacketBuilder().
			DestinationIP(icmpPacket.Source()).
			SourceIP(icmpPacket.Destination()).
			TTL(10).
			Protocol(ICMPProtocol).
			Build()
		echo := &icmp{
			ipEcho,
			ICMPEchoReply,
		}
		fIPSend(node, echo)
	}
}
