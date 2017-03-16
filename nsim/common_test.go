package nsim

import (
	"net"
	"testing"
)

type dummyMedium struct {
}

func (dm *dummyMedium) Send(frame Frame) error {
	return nil
}

func NewMockNode(t *testing.T) *MockNode {
	return &MockNode{t, nil, nil, nil}
}

type MockNode struct {
	T                  *testing.T
	FNetworkInterfaces func() []NetworkInterface
	FMedium            func() TransmissionMedium
	FRoutingTable      func() []Route
}

func (mn *MockNode) NetworkInterfaces() []NetworkInterface {
	if mn.FNetworkInterfaces == nil {
		mn.T.Error("NetworkInterfaces() handler is not set in mock.")
		return nil
	}
	return mn.FNetworkInterfaces()
}

func (mn *MockNode) Medium() TransmissionMedium {
	if mn.FMedium == nil {
		mn.T.Error("Medium() handler is not set in mock.")
		return nil
	}
	return mn.FMedium()
}

func (mn *MockNode) RoutingTable() []Route {
	if mn.FRoutingTable == nil {
		mn.T.Error("RoutingTable() handler is not set in mock.")
		return nil
	}
	return mn.FRoutingTable()
}

func NewMockIPPacket(t *testing.T) *MockIPPacket {
	return &MockIPPacket{t, nil, nil, nil, nil, nil}
}

type MockIPPacket struct {
	T            *testing.T
	FDestination func() net.IP
	FSource      func() net.IP
	FTTL         func() int
	FProtocol    func() string
	FDecreaseTTL func()
}

func (mip *MockIPPacket) Destination() net.IP {
	if mip.FDestination == nil {
		mip.T.Error("Destination() handler is not set in mock.")
		return nil
	}
	return mip.FDestination()
}

func (mip *MockIPPacket) Source() net.IP {
	if mip.FSource == nil {
		mip.T.Error("Source() handler is not set in mock.")
		return nil
	}
	return mip.FSource()
}

func (mip *MockIPPacket) TTL() int {
	if mip.FTTL == nil {
		mip.T.Error("TTL() handler is not set in mock.")
		return -1
	}
	return mip.FTTL()
}

func (mip *MockIPPacket) Protocol() string {
	if mip.FProtocol == nil {
		mip.T.Error("Protocol() handler is not set in mock.")
		return ""
	}
	return mip.FProtocol()
}

func (mip *MockIPPacket) DecreaseTTL() {
	if mip.FDecreaseTTL == nil {
		mip.T.Error("DecreaseTTL is not implemented.")
	}
	mip.FDecreaseTTL()
}
