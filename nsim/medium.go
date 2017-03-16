package nsim

// SingleQueueMedium is simple one threaded frame delivery based on queue.
type SingleQueueMedium struct {
	nodes map[string]Node
	head  *queueItem
	tail  *queueItem
}

// NewSignleQueueMedium creates new medium.
func NewSignleQueueMedium() *SingleQueueMedium {
	return &SingleQueueMedium{make(map[string]Node), nil, nil}
}

type queueItem struct {
	next  *queueItem
	frame *Frame
}

// RegisterNode registers a node for delivery. Panics if an ip belongs to more
// than one node.
func (m *SingleQueueMedium) RegisterNode(node Node) {
	for _, ni := range node.NetworkInterfaces() {
		mac := ni.IP.String()
		if _, exists := m.nodes[mac]; exists {
			panic("IP is already registered.")
		}
		m.nodes[mac] = node
	}
}

// Send schedules a frame for delivery.
func (m *SingleQueueMedium) Send(frame Frame) error {
	n := &queueItem{nil, &frame}
	if m.head == nil {
		m.head = n
		m.tail = n
	} else {
		m.tail.next = n
		m.tail = n
	}
	return nil
}

// DeliverFrame processes next frame if there is one. Does nothing if no frames
// to process. Returns true is queue is not empty after one frame was consumed.
func (m *SingleQueueMedium) DeliverFrame() bool {
	if m.head == nil {
		return false
	}
	queueItem := m.head
	m.head = m.head.next
	if m.head == nil {
		m.tail = nil
	}
	m.processFrame(queueItem.frame)
	return m.head != nil
}

func (m *SingleQueueMedium) processFrame(frame *Frame) {
	node, exists := m.nodes[frame.destinationID]
	if !exists {
		return
	}
	fLinkReceive(node, *frame)
}

// HasMore returns true if frames queue is not empty.
func (m *SingleQueueMedium) HasMore() bool {
	return m.head != nil
}
