package sequencefile

import (
	"errors"
	"testing"
)

func TestBinarySearchIndex(t *testing.T) {
	var tests = []struct {
		input     []*Index
		index     int64
		block     int64
		blockNext int64
		err       error
	}{
		{[]*Index{}, 1, -1, -1, ErrEmptyIndexs},
		{[]*Index{{Index: 1}}, 0, -1, -1, ErrIndexNotFound},
		{[]*Index{{Index: 1}}, 1, 1, -1, nil},
		{[]*Index{{Index: 1}}, 2, 1, -1, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}}, 1, 1, 5, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}}, 11, 10, 15, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}}, 50, 20, -1, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}, {Index: 30}}, 10, 10, 15, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}, {Index: 30}}, 11, 10, 15, nil},
		{[]*Index{{Index: 1}, {Index: 5}, {Index: 10}, {Index: 15}, {Index: 20}, {Index: 30}}, 50, 30, -1, nil},
	}
	for i, test := range tests {
		b, bn, err := binarySearchIndex(test.input, test.index)
		if err != nil {
			if test.err == nil {
				t.Errorf("test[%d]: expect nil err, got err: %v", i, err)
				continue
			}
			if errors.Is(err, test.err) {
				continue
			}
			t.Errorf("test[%d]: expect err: %v, got err: %v", i, test.err, err)
			continue
		}
		if b.Index != test.block {
			t.Errorf("test[%d]: expect block %d, got block %d", i, test.block, b.Index)
			continue
		}
		if test.blockNext == -1 {
			if bn != nil {
				t.Errorf("test[%d]: expect blockNext nil, got blockNext %d", i, bn.Index)
			}
			continue
		}
		if test.blockNext != bn.Index {
			t.Errorf("test[%d]: expect blockNext %d, got blockNext %d", i, test.blockNext, bn.Index)
		}
	}
}
