// Link layer of TCP/IP stack (App-Transport-Network-Link).
// This simulation is not using MAC addresses. Nodes have IP addresses assigned
// to them; these IP addresses are used instead of MACs to simplify configuration.

package nsim

import "net"

// Frame is a model to simulate link layer in TCP/IP stack.
type Frame struct {
	destinationID string
}

// LinkReceive is called when a node has an incoming frame. The receiver may
// ignore the frame if it's not a target. This behaviour simulates ethernet network.
// See package header comment for more info about MAC/IP addresses.
func (node *Node) LinkReceive(frame Frame) bool {
	targetIP := net.ParseIP(frame.destinationID)
	for _, ni := range node.NetworkInterfaces {
		if ni.IP.Equal(targetIP) {
			return true
		}
	}
	return false
}
