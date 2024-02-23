package invcount

import (
	"slices"
	"testing"
)

func TestCountInversions(t *testing.T) {
	cases := []struct {
		input    []int
		expected int
	}{
		{[]int{1, 2, 3, 4, 5}, 0},    // Sorted array
		{[]int{2, 4, 1, 3, 5}, 3},    // (2, 1), (4, 1), (4, 3)
		{[]int{5, 4, 3, 2, 1}, 10},   // Maximum inversions for a reversed array
		{[]int{3, 5, 1, 2, 4}, 5},    // (3, 1), (3, 2), (5, 1), (5, 2), (5,4)
		{[]int{1, 1, 1, 1, 1}, 0},    // All elements are same, no inversions
		{[]int{2, 1}, 1},             // (2, 1)
		{[]int{}, 0},                 // Empty array
		{[]int{1}, 0},                // Single element array
		{[]int{1, 3, 5, 2, 4, 6}, 3}, // (3, 2), (5, 2), (5, 4)
		{[]int{5, 3, 1, 2, 4}, 6},    // (5, 3), (5, 1), (5, 2), (5, 4), (3, 1), (3, 2)
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0},
		{[]int{9, 8, 7, 6, 5, 4, 3, 2, 1}, 36},
		{[]int{5, 3, 7, 2, 8, 4, 6, 1}, 16},
		{[]int{1, 4, 2, 5, 7, 4, 2, 4, 6, 0, 7, 4, 6, 7, 3, 2, 5, 6, 4, 7, 3, 6, 1}, 100},
	}

	for _, c := range cases {
		cc := append([]int(nil), c.input...)
		got := Count(cc)
		if got != c.expected {
			t.Errorf("Count(%v) == %d, expected %d", c.input, got, c.expected)
		}
		if !slices.Equal(cc, c.input) {
			t.Errorf("Count(%v) modified input", c.input)
		}
	}
}
