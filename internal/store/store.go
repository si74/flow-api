package store

import (
	"container/list"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flow represents a single data point of network bytes transmitted and received for a given tuple
// A flow may be a single data point or an aggregation of data points.
type Flow struct {
	// Src is the source application name
	Src string `json:"src_app"`
	// Dst is the destination application name
	Dst     string `json:"dst_app"`
	VpcID   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytsRx  int    `json:"bytes_rx"`
	Hour    int    `json:"hour"`
}

// FlowKey represents a unique tuple of identifying flow characteristics
type FlowKey struct {
	Src   string
	Dst   string
	VpcID string
}

// Could use a sync map but wanted to try this instead
// TOOD(sneha) Thread-safe but there can be lock contention
// How to speed this up?
// FlowStore is a mapping of a linked list of flow data - keyed by a unique flow tuple
type FlowStore struct {
	mu sync.RWMutex
	// a mapping of a uniquely identifying flow key to a doubly-linked list of flows
	flowMap map[FlowKey]*FlowList
	// TODO(Sneha): make this configurable
	//maxSize int
	// retention
	// TODO(sneha): add some metrics
	// lastUpdated
	ll *logrus.Logger
}

// NewFlowStore creates and returns a new flow store
func NewFlowStore() *FlowStore {
	return &FlowStore{
		mu:      sync.RWMutex{},
		flowMap: map[FlowKey]*FlowList{},
	}
}

// Insert adds a flow entry in chronological order for a flowList
func (fs *FlowStore) Insert(flows []*Flow) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	for _, flow := range flows {
		key := FlowKey{
			Src:   flow.Src,
			Dst:   flow.Dst,
			VpcID: flow.VpcID,
		}
		flowList, ok := fs.flowMap[key]
		if !ok {
			flowList = &FlowList{l: list.New()}
			fs.flowMap[key] = flowList
		}
		err := flowList.insert(flow)
		if err != nil {
			// Log and continue
		}
	}

	return nil
}

type FlowList struct {
	l *list.List
}

// Insert inserts a flow for a current hour or sums values into existing flow entry
// Note: For this case, it was unnecessary to store all input flow points separately if
// they are for the same flow key and time window. In the future, this may change.
func (fl *FlowList) insert(flow *Flow) error {

	return nil
}

// Get returns an aggregation of flow stats for all tuples for a given hour
func (fs *FlowStore) Get(hour int) ([]*Flow, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return nil, nil
}
