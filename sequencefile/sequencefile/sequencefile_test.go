package sequencefile

import (
	"errors"
	"os"
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

func benchmarkWriteLength(b *testing.B, msgLength int) {
	var s SeqFile
	os.RemoveAll("data")
	defer os.RemoveAll("data")
	s.Init()
	defer s.Close()
	msg := RandBytes(int64(msgLength), msgLength)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Write(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

// BenchmarkWrite100 writes messages with a size of 100 bytes
func BenchmarkWrite100(b *testing.B) {
	benchmarkWriteLength(b, 100)
}

// BenchmarkWrite1k writes messages with a size of 1000 bytes
func BenchmarkWrite1k(b *testing.B) {
	benchmarkWriteLength(b, 1000)
}

// BenchmarkWrite10k writes messages with a size of 10,000 bytes (10kb)
func BenchmarkWrite10k(b *testing.B) {
	benchmarkWriteLength(b, 10000)
}

// BenchmarkWrite100k writes messages with a size of 100,000 bytes (100kb)
func BenchmarkWrite100k(b *testing.B) {
	benchmarkWriteLength(b, 10000)
}

func benchmarkParallelWriteLength(b *testing.B, msgLength int) {
	var s SeqFile
	os.RemoveAll("data")
	defer os.RemoveAll("data")
	s.Init()
	defer s.Close()
	msg := RandBytes(int64(msgLength), msgLength)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := s.Write(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
}

func BenchmarkParallelWrite100(b *testing.B) {
	benchmarkParallelWriteLength(b, 100)
}

func BenchmarkParallelWrite1k(b *testing.B) {
	benchmarkParallelWriteLength(b, 1000)
}

func BenchmarkParallelWrite10k(b *testing.B) {
	benchmarkParallelWriteLength(b, 10000)
}

func BenchmarkParallelWrite100k(b *testing.B) {
	benchmarkParallelWriteLength(b, 100000)
}

