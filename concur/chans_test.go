package concur

import (
	"context"
	"testing"
)

func toChan(ss []int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range ss {
			c <- v
		}
		close(c)
	}()
	return c
}

func TestFanIn(t *testing.T) {
	tests := []struct {
		cnt   int
		parts []int
	}{
		{cnt: 3, parts: nil},
		{cnt: 100, parts: nil},
		{cnt: 3, parts: []int{1}},
		{cnt: 3, parts: []int{2}},
		{cnt: 3, parts: []int{1, 2}},
		{cnt: 10, parts: []int{1}},
		{cnt: 10, parts: []int{5}},
		{cnt: 10, parts: []int{1, 9}},
		{cnt: 3, parts: []int{0, 0, 0, 2}},
		{cnt: 10, parts: []int{1, 2, 3, 8}},
		{cnt: 10, parts: []int{1, 2, 3, 4, 4, 8}},
		{cnt: 1024, parts: []int{300, 1000}},
		{cnt: 1024, parts: []int{20, 21, 30, 50, 100, 134, 500, 789, 1000}},
		{cnt: 1024, parts: []int{0, 0, 20, 21, 30, 50, 100, 134, 134, 134, 500, 501, 502, 502, 600, 600, 789, 1000}},
	}

	for _, tt := range tests {
		var ss []int
		for v := range tt.cnt {
			ss = append(ss, v)
		}

		var chans []<-chan int
		i := 0
		for _, v := range tt.parts {
			chans = append(chans, toChan(ss[i:v]))
			i = v
		}
		chans = append(chans, toChan(ss[i:]))

		mergedCh := FanIn(context.Background(), chans...)

		dedup := make(map[int]struct{})
		for v := range mergedCh {
			if _, ok := dedup[v]; ok {
				t.Errorf("duplicated value: %d", v)
				continue
			}
			dedup[v] = struct{}{}
		}
		if len(dedup) != tt.cnt {
			t.Errorf("FanIn() got: %d unique values, want: %d", len(dedup), tt.cnt)
		}
	}
}
