package integrationtests

import (
	"github.com/andrewromanenco/network-sim/nsim"
	"testing"
)

func TestPing(t *testing.T) {
	internet, nodeFrom, nodeTo := createNetwork()
	pingRequest := createPingRequest(nodeFrom, nodeTo)

	nsim.IPSend(nodeFrom, pingRequest)

	if !internet.HasMore() {
		t.Error("Ping packet must be available for processing")
	}
	count := 1
	var lastFrame *nsim.Frame
	for internet.DeliverFrame() {
		count++
		lastFrame = internet.Frame()
	}
	if count != 4 {
		t.Error("Expected 4 frames to be delivered. From->Router, R->To, T->R, R->F.")
	}
	if lastFrame == nil {
		t.Error("Should never happen.")
	}
	icmpEcho, ok := lastFrame.IPPacket.(nsim.ICMPPacket)
	if !ok {
		t.Error("Last frame expected to be icmp packet.")
		return
	}
	if !icmpEcho.Destination().Equal(nodeFrom.NetworkInterfaces()[0].IP) {
		t.Error("Last packet should be for node1.")
	}
	if !icmpEcho.Source().Equal(nodeTo.NetworkInterfaces()[0].IP) {
		t.Error("Last packet should be from node2.")
	}
	if icmpEcho.MessageType() != nsim.ICMPEchoReply {
		t.Error("Last packet must be ping echo reply.")
	}
	if icmpEcho.TTL() != 9 {
		t.Error("TTL expected to be reduced to 9 as packet was routed onece.")
	}
}

func createNetwork() (*nsim.SingleQueueMedium, nsim.Node, nsim.Node) {
	nsim.RegisterProtocolHandler(nsim.ICMPProtocol, nsim.ICMPHandler)
	medium := nsim.NewSignleQueueMedium()
	node1 := nsim.NewNodeBuilder().
		AddNetInterface("192.168.1.10/24").
		AddRoute("0.0.0.0/0", "192.168.1.1").
		WithMedium(medium).
		Build()
	node2 := nsim.NewNodeBuilder().
		AddNetInterface("192.168.2.10/24").
		AddRoute("0.0.0.0/0", "192.168.2.1").
		WithMedium(medium).
		Build()
	router := nsim.NewNodeBuilder().
		AddNetInterface("192.168.1.1/24").
		AddNetInterface("192.168.2.1/24").
		AddRoute("0.0.0.0/0", "0.0.0.0").
		WithMedium(medium).
		Build()
	medium.RegisterNode(node1)
	medium.RegisterNode(node2)
	medium.RegisterNode(router)
	return medium, node1, node2
}

func createPingRequest(nodeFrom nsim.Node, nodeTo nsim.Node) nsim.ICMPPacket {
	ipPacket := nsim.NewIPPacketBuilder().
		DestinationIP(nodeTo.NetworkInterfaces()[0].IP).
		SourceIP(nodeFrom.NetworkInterfaces()[0].IP).
		TTL(10).
		Protocol(nsim.ICMPProtocol).
		Build()
	return nsim.NewICMPPacket(
		ipPacket,
		nsim.ICMPEchoRequest,
	)
}
