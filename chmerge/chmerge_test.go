package chmerge

import (
	"fmt"
	"testing"

	stdcmp "cmp"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		chCount int
		elCount int
	}{
		{chCount: 0, elCount: 0},
		{chCount: 1, elCount: 0},
		{chCount: 1, elCount: 1},
		{chCount: 1, elCount: 10},
		{chCount: 2, elCount: 0},
		{chCount: 2, elCount: 1},
		{chCount: 2, elCount: 2},
		{chCount: 2, elCount: 10},
		{chCount: 3, elCount: 0},
		{chCount: 3, elCount: 1},
		{chCount: 3, elCount: 2},
		{chCount: 3, elCount: 10},
		{chCount: 10, elCount: 0},
		{chCount: 10, elCount: 1},
		{chCount: 10, elCount: 10},
		{chCount: 10, elCount: 20},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("channels: %d, elements: %d", tc.chCount, tc.elCount), func(t *testing.T) {
			cc := createChannels(tc.chCount)
			writeToChannels(tc.elCount, cc...)

			mergedCh := Merge(FromSlice(cc)...)

			// Receive values from the merged channel
			var result []int
			for range tc.elCount * tc.chCount {
				result = append(result, <-mergedCh)
			}

			if _, ok := <-mergedCh; ok {
				t.Errorf("expected merged channel to be closed")
			}

			// Compare the received values with the expected values
			var expected []int
			for i := range tc.elCount * tc.chCount {
				expected = append(expected, i)
			}

			if !cmp.Equal(expected, result, cmpopts.SortSlices(stdcmp.Less[int])) {
				t.Fatalf("expected: %v; got: %v", expected, result)
			}
		})
	}
}

func createChannels(n int) []chan int {
	var cc []chan int
	for range n {
		cc = append(cc, make(chan int))
	}
	return cc
}

func writeToChannels(n int, channels ...chan int) {
	for i, c := range channels {
		c := c
		go func(beg int) {
			for i := range n {
				c <- beg*n + i
			}
			close(c)
		}(i)
	}
}
