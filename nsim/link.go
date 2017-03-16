// Link layer of TCP/IP stack (App-Transport-Network-Link).
// This simulation is not using MAC addresses. Nodes have IP addresses assigned
// to them; these IP addresses are used instead of MACs to simplify configuration.

package nsim

import "net"

// TransmissionMedium is an abstraction for sending frames between nodes.
type TransmissionMedium interface {
	Send(frame Frame) error
}

// Frame is a model to simulate link layer in TCP/IP stack.
type Frame struct {
	destinationID string
	IPPacket      IPPacket
}

var fIPReceive = IPReceive
var fLinkReceive = LinkReceive

// LinkReceive is called when a node has an incoming frame. The receiver may
// ignore the frame if it's not a target. This behaviour simulates ethernet network.
// See package header comment for more info about MAC/IP addresses.
func LinkReceive(node Node, frame Frame) bool {
	targetIP := net.ParseIP(frame.destinationID)
	for _, ni := range node.NetworkInterfaces() {
		if ni.IP.Equal(targetIP) {
			fIPReceive(node, frame.IPPacket)
			return true
		}
	}
	return false
}

// LinkSend sends a frame to the transmission medium.
func LinkSend(node Node, frame Frame) error {
	return node.Medium().Send(frame)
}
