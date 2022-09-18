package store

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

// less if used for sorting Flow slices for testing
func less(f1 *Flow, f2 *Flow) bool {
	if f1.Src != f2.Src {
		return f1.Src < f2.Src
	}
	if f1.Dst != f2.Dst {
		return f1.Dst < f2.Dst
	}
	if f1.VpcID != f2.VpcID {
		return f1.VpcID < f2.VpcID
	}
	if f1.BytesTx != f2.BytesTx {
		return f1.BytesTx < f2.BytesTx
	}
	if f1.BytesRx != f2.BytesRx {
		return f1.BytesRx < f2.BytesRx
	}
	return f1.Hour < f2.Hour
}

func Test_Flow(t *testing.T) {
	// TODO(sneha): Add more extensive tests in table tests later
	tests := []struct {
		name          string
		insert        []*Flow
		Hour          int
		expectedFlows []*Flow
	}{
		{
			name:          "empty flow database",
			insert:        nil,
			Hour:          1,
			expectedFlows: []*Flow{},
		},
		{
			name: "0 flows returned",
			insert: []*Flow{
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 100, BytesRx: 300, Hour: 1},
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 200, BytesRx: 600, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 2},
				{Src: "baz", Dst: "qux", VpcID: "vpc-1", BytesTx: 100, BytesRx: 500, Hour: 2},
			},
			Hour:          3,
			expectedFlows: []*Flow{},
		},
		{
			name: "multiple flows returned",
			insert: []*Flow{
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 100, BytesRx: 300, Hour: 1},
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 200, BytesRx: 600, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 2},
				{Src: "baz", Dst: "qux", VpcID: "vpc-1", BytesTx: 100, BytesRx: 500, Hour: 2},
			},
			Hour: 2,
			expectedFlows: []*Flow{
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 2},
				{Src: "baz", Dst: "qux", VpcID: "vpc-1", BytesTx: 100, BytesRx: 500, Hour: 2},
			},
		},
		{
			name: "multiple aggregated flows returned",
			insert: []*Flow{
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 100, BytesRx: 300, Hour: 1},
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 200, BytesRx: 600, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 2},
				{Src: "baz", Dst: "qux", VpcID: "vpc-1", BytesTx: 100, BytesRx: 500, Hour: 2},
			},
			Hour: 1,
			expectedFlows: []*Flow{
				{Src: "foo", Dst: "bar", VpcID: "vpc-0", BytesTx: 300, BytesRx: 900, Hour: 1},
				{Src: "baz", Dst: "qux", VpcID: "vpc-0", BytesTx: 100, BytesRx: 500, Hour: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ll := logrus.New()
			ll.SetOutput(io.Discard)

			store := NewFlowStore(ll)

			if tt.insert != nil {
				err := store.Insert(tt.insert)
				if err != nil {
					t.Fatalf("unexpected error inserting flows: %v", err)
				}
			}

			flows, err := store.Get(tt.Hour)
			if err != nil {
				t.Fatalf("unexpected error retrieving flows: %v", err)
			}

			if diff := cmp.Diff(flows, tt.expectedFlows, cmpopts.SortSlices(less)); diff != "" {
				t.Fatalf("unexpected flows: %v", diff)
			}
		})
	}
}
