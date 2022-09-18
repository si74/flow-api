package store

import (
	"container/list"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flow represents a single data point of network bytes transmitted and received for a given tuple
// A flow may be a single data point or an aggregation of data points.
type Flow struct {
	// Src is the source application name
	Src string `json:"src_app"`
	// Dst is the destination application name
	Dst     string `json:"dest_app"`
	VpcID   string `json:"vpc_id"`
	BytesTx int    `json:"bytes_tx"`
	BytesRx int    `json:"bytes_rx"`
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
	// a mapping of a uniquely identifying flow key to a generic flow list interface
	flowMap map[FlowKey]flowList
	// TODO(Sneha): make this configurable
	// maxSize int
	// retention int
	// TODO(sneha): add some metrics
	// lastUpdated
	ll *logrus.Logger
}

// NewFlowStore creates and returns a new flow store
func NewFlowStore(ll *logrus.Logger) *FlowStore {
	return &FlowStore{
		mu:      sync.RWMutex{},
		flowMap: map[FlowKey]flowList{},
		ll:      ll,
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
			flowList = &flowListV1{l: list.New()}
			fs.flowMap[key] = flowList
		}
		err := flowList.insert(flow)
		if err != nil {
			fs.ll.Errorf("unable to insert flow for %v: %v", key, err)
			continue
		}
	}
	return nil
}

// Get returns an aggregation of flow stats for all tuples for a given hour
func (fs *FlowStore) Get(hour int) ([]*Flow, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if hour <= 0 {
		return nil, errors.New("timestamp must be greater than 0")
	}

	flows := []*Flow{}
	for key, list := range fs.flowMap {
		flow, err := list.get(key, hour)
		if err != nil {
			fs.ll.Errorf("unable to retrieve aggregate flow for %v: %v", key, err)
			continue
		}
		// No data present
		if flow == nil {
			continue
		}
		flows = append(flows, flow)
	}
	return flows, nil
}

// flowList is a generic interface that accepts flow data points
// returns aggregated flow data.
// Implementations may vary in run-time complexity and efficiency.
type flowList interface {
	insert(flow *Flow) error
	get(key FlowKey, hour int) (*Flow, error)
}

// flowListV1 is the less optimized flow list that inserts in any order and does a brute-force
// aggregation of flow values.
type flowListV1 struct {
	l *list.List
}

// insert a new flow data point.
// Note: The current iteration does not preserve chronological order, which makes insertion quite fast.
// To improve aggregation speed in the future, chronological insertion or aggregation into a single flow on insert may be done.
func (fl *flowListV1) insert(flow *Flow) error {
	// Case when the list is currently empty and we can insert the flow
	fl.l.PushBack(flow)
	return nil
}

// get returns an aggregated for for a given tuple and hour timestamp
func (fl *flowListV1) get(key FlowKey, hour int) (*Flow, error) {
	if hour <= 0 {
		return nil, fmt.Errorf("provided hour timestamp must be greater than 0")
	}

	var found bool
	aggregateFlow := &Flow{
		Src:     key.Src,
		Dst:     key.Dst,
		VpcID:   key.VpcID,
		Hour:    hour,
		BytesTx: 0,
		BytesRx: 0,
	}

	for e := fl.l.Front(); e != nil; e = e.Next() {
		flow, ok := e.Value.(*Flow)
		if !ok {
			e = e.Next()
			continue
		}

		if flow.Hour == hour {
			found = true
			aggregateFlow.BytesRx += flow.BytesRx
			aggregateFlow.BytesTx += flow.BytesTx
		}
	}

	// This timestamp was never found, return a nil flow value
	if !found {
		return nil, nil
	}

	return aggregateFlow, nil
}
